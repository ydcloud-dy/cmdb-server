package configCenter

import (
	"DYCLOUD/global"
	"DYCLOUD/model/configCenter"
	"DYCLOUD/model/configCenter/request"
	cicd2 "DYCLOUD/utils/cicd"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-atomci/workflow/jenkins"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ServiceIntegrationService struct{}

// GetServiceIntegrationList
//
//	@Description: 获取服务列表
//	@receiver s
//	@param req
//	@return data
//	@return total
//	@return err
func (s *ServiceIntegrationService) GetServiceIntegrationList(req *request.ServiceRequest) (data *[]configCenter.ServiceIntegration, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	// 创建db
	db := global.DYCLOUD_DB.Model(&configCenter.ServiceIntegration{})
	var serviceList []configCenter.ServiceIntegration
	// 如果有条件搜索 下方会自动创建搜索语句
	//if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
	//db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt).Where("project = ?",info.Project)
	//db = db.Where("name = ?", req.Keyword)
	//}
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		db.Where("name like ?", keyword).Or("id = ?", req.Keyword)
	}
	db.Where("type != 3")
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&serviceList).Error
	if err != nil {
		return nil, 0, err
	}

	// 对每个配置进行解密并存储到 `DecryptedConfig`
	for i := range serviceList {
		if err := serviceList[i].DecryptConfig(); err != nil {
			fmt.Printf("Error decrypting config for service %s: %v\n", serviceList[i].Name, err)
		}
	}
	return &serviceList, total, err
}
func (s *ServiceIntegrationService) GetRegistryList(req *request.ServiceRequest) (data *[]configCenter.ServiceIntegration, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	// 创建db
	db := global.DYCLOUD_DB.Model(&configCenter.ServiceIntegration{})
	var serviceList []configCenter.ServiceIntegration

	db.Where("type = 1")
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&serviceList).Error
	if err != nil {
		return nil, 0, err
	}

	// 对每个配置进行解密并存储到 `DecryptedConfig`
	for i := range serviceList {
		if err := serviceList[i].DecryptConfig(); err != nil {
			fmt.Printf("Error decrypting config for service %s: %v\n", serviceList[i].Name, err)
		}
	}
	return &serviceList, total, err
}

// CreateServiceIntegration
//
//	@Description: 创建服务
//	@receiver s
//	@param req
//	@return error
func (s *ServiceIntegrationService) CreateServiceIntegration(req *configCenter.ServiceIntegration) error {
	config, err := req.ChooseConfig()
	if err != nil {
		return err
	}

	req.CryptoConfig(config)
	var existingService configCenter.ServiceIntegration
	if err := global.DYCLOUD_DB.Where("name = ?", req.Name).First(&existingService).Error; err == nil {
		return fmt.Errorf("服务名称 '%s' 已存在，请选择一个唯一的名称", req.Name)
	}
	if err := global.DYCLOUD_DB.Model(&configCenter.ServiceIntegration{}).Create(&req).Error; err != nil {
		return err
	}
	return nil
}

// UpdateServiceIntegration
//
//	@Description: 更新服务
//	@receiver s
//	@param req
//	@return error
func (s *ServiceIntegrationService) UpdateServiceIntegration(req *configCenter.ServiceIntegration) error {
	data, err := s.DescribeServiceIntegration(int(req.ID))
	if err != nil {
		return err
	}
	data = req
	config, err := req.ChooseConfig()
	if err != nil {
		return err
	}
	data.CryptoConfig(config)
	if err = global.DYCLOUD_DB.Model(&configCenter.ServiceIntegration{}).Where("id = ?", req.ID).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

// DeleteServiceIntegration
//
//	@Description: 删除服务
//	@receiver s
//	@param id
//	@return error
func (s *ServiceIntegrationService) DeleteServiceIntegration(id int) error {

	if err := global.DYCLOUD_DB.Where("id = ?", id).Delete(&configCenter.ServiceIntegration{}).Error; err != nil {
		return err
	}
	return nil
}

// DescribeServiceIntegration
//
//	@Description: 获取服务详情
//	@receiver s
//	@param id
//	@return *configCenter.ServiceIntegration
//	@return error
func (s *ServiceIntegrationService) DescribeServiceIntegration(id int) (*configCenter.ServiceIntegration, error) {
	var data configCenter.ServiceIntegration
	if err := global.DYCLOUD_DB.Model(&configCenter.ServiceIntegration{}).Where("id = ? and type != 3", id).First(&data).Error; err != nil {
		return nil, err
	}
	data.DecryptConfig()
	return &data, nil
}

// VerifyServiceIntegration
//
//	@Description: 验证服务是否可以连接
//	@receiver s
//	@param req
//	@return string
//	@return error
func (s *ServiceIntegrationService) VerifyServiceIntegration(req *configCenter.ServiceIntegration) (string, error) {
	switch req.Type {
	case configCenter.KUBERNETES_TYPE:
		kube := &configCenter.KubernetesConfig{}
		err := json.Unmarshal([]byte(req.Config), kube)
		if err != nil {
			global.DYCLOUD_LOG.Error(fmt.Sprintf("kubernetes conf format error: %v", err.Error()))
			return "", err
		}
		var k8sConf *rest.Config
		switch kube.Type {
		case configCenter.KUBERNETES_CONFIG:
			k8sConf, err = clientcmd.RESTConfigFromKubeConfig([]byte(kube.Conf))
			if err != nil {
				return "", err
			}
		case configCenter.KUBERNETES_TOKEN:
			k8sConf = &rest.Config{
				BearerToken:     kube.Conf,
				TLSClientConfig: rest.TLSClientConfig{Insecure: true},
				Host:            kube.Url,
			}

		}
		clientset, err := kubernetes.NewForConfig(k8sConf)
		if err != nil {
			return "", err
		}
		k8sVersion, err := clientset.Discovery().ServerVersion()
		if err != nil {
			global.DYCLOUD_LOG.Error(fmt.Sprintf("获取k8s版本失败：%s", err.Error()))
			return "", err
		}
		msg := fmt.Sprintf("成功连接k8s集群： %s", k8sVersion.GitVersion)
		return msg, nil

	case configCenter.REGISTRY_TYPE:
		registryConf := &configCenter.RegistryConfig{}
		err := json.Unmarshal([]byte(req.Config), registryConf)
		if err != nil {
			global.DYCLOUD_LOG.Error(fmt.Sprintf("registryConf unmarshal failed", err.Error()))
			return "", err
		}
		err = cicd2.TryLoginRegistry(registryConf.Url, registryConf.Username, registryConf.Password, registryConf.IsHttps)
		if err != nil {
			global.DYCLOUD_LOG.Error(fmt.Sprintf("登录仓库失败:%s", err.Error()))
			return "", err
		}
		return "连接成功", err
	case configCenter.JENKINS_TYPE:
		jenkinsConf := &configCenter.JenkinsConfig{}
		err := json.Unmarshal([]byte(req.Config), jenkinsConf)
		if err != nil {
			global.DYCLOUD_LOG.Error(fmt.Sprintf("jenkinsConf unmarshal failed", err.Error()))
			return "", err
		}
		global.DYCLOUD_LOG.Debug(fmt.Sprintf("verify jenkins conf:%v", jenkinsConf))
		jClient, err := jenkins.NewJenkinsClient(
			jenkins.URL(jenkinsConf.Url),
			jenkins.JenkinsUser(jenkinsConf.Username),
			jenkins.JenkinsToken(jenkinsConf.Token),
		)
		if err != nil {
			global.DYCLOUD_LOG.Error(fmt.Sprintf("创建jenkins客户端失败:%s", err.Error()))
			return "", err
		}
		pingInfo, err := jClient.Ping()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("成功连接jenkins：%s", pingInfo), nil
	default:
		return "", errors.New(fmt.Sprintf("不支持的类型：%s", req.Type))
	}
	return "", nil

}
