package cicd

import (
	"DYCLOUD/global"
	"DYCLOUD/model/cicd"
	"DYCLOUD/model/cicd/request"
	"DYCLOUD/service/cicd/dao"
	"encoding/json"
	"fmt"
)

type SourceCodeService struct{}

func (s *SourceCodeService) GetSourceCodeList(req *request.ServiceRequest) (data *[]cicd.ServiceIntegration, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	// 创建db
	db := global.DYCLOUD_DB.Model(&cicd.ServiceIntegration{})
	var serviceList []cicd.ServiceIntegration
	// 如果有条件搜索 下方会自动创建搜索语句
	//if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
	//db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt).Where("project = ?",info.Project)
	//db = db.Where("name = ?", req.Keyword)
	//}
	db.Where("type = 3")
	if req.Keyword != "" {
		db.Where("name = ?", req.Keyword).Or("id = ?", req.Keyword)
	}
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

// CreateSourceCode
//
//	@Description: 创建代码源
//	@receiver s
//	@param req
//	@return error
func (s *SourceCodeService) CreateSourceCode(req *cicd.ServiceIntegration) error {
	config, err := req.ChooseConfig()
	if err != nil {
		return err
	}

	req.CryptoConfig(config)

	if err := global.DYCLOUD_DB.Model(&cicd.ServiceIntegration{}).Create(&req).Error; err != nil {
		return err
	}
	return nil
}

// UpdateSourceCode
//
//	@Description: 更新代码源
//	@receiver s
//	@param req
//	@return error
func (s *SourceCodeService) UpdateSourceCode(req *cicd.ServiceIntegration) error {
	fmt.Println(req)
	data, err := s.DescribeSourceCode(int(req.ID))
	if err != nil {
		return err
	}
	data = req
	config, err := req.ChooseConfig()
	if err != nil {
		return err
	}
	data.CryptoConfig(config)
	if err = global.DYCLOUD_DB.Model(&cicd.ServiceIntegration{}).Where("id = ?", req.ID).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

// DeleteSourceCode
//
//	@Description: 删除代码源
//	@receiver s
//	@param id
//	@return error
func (s *SourceCodeService) DeleteSourceCode(id int) error {
	fmt.Println(id)

	if err := global.DYCLOUD_DB.Where("id = ?", id).Delete(&cicd.ServiceIntegration{}).Error; err != nil {
		return err
	}
	return nil
}

// DescribeSourceCode
//
//	@Description: 查看代码源详情
//	@receiver s
//	@param id
//	@return *cicd.ServiceIntegration
//	@return error
func (s *SourceCodeService) DescribeSourceCode(id int) (*cicd.ServiceIntegration, error) {
	fmt.Println(id)
	var data cicd.ServiceIntegration
	if err := global.DYCLOUD_DB.Model(&cicd.ServiceIntegration{}).Where("id = ? and type = 3", id).First(&data).Error; err != nil {
		return nil, err
	}
	data.DecryptConfig()
	return &data, nil
}

// VerifySourceCode
//
//	@Description: 验证服务是否可以连接
//	@receiver s
//	@param req
//	@return string
//	@return error
func (s *SourceCodeService) VerifySourceCode(req *cicd.ServiceIntegration) (string, error) {
	fmt.Println(req)
	gitConf := &cicd.GitConfig{}
	err := json.Unmarshal([]byte(req.Config), gitConf)
	if err != nil {
		return "", err
	}
	err = dao.VerifyRepoConnetion(req.Type, gitConf.Url, gitConf.Token)
	if err != nil {
		return "", err
	}
	return "连接成功", nil

}
