package configCenter

import (
	"DYCLOUD/global"
	"DYCLOUD/model/configCenter"
	"DYCLOUD/model/configCenter/request"
	"DYCLOUD/service/configCenter/dao"
	"encoding/json"
	"fmt"
	"github.com/drone/go-scm/scm"
	"golang.org/x/net/context"
)

type SourceCodeService struct{}

func (s *SourceCodeService) GetSourceCodeList(req *request.ServiceRequest) (data *[]configCenter.ServiceIntegration, total int64, err error) {
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
	db.Where("type in (3,4,5,6,7)")
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
func (s *SourceCodeService) CreateSourceCode(req *configCenter.ServiceIntegration) error {
	config, err := req.ChooseConfig()
	if err != nil {
		return err
	}

	req.CryptoConfig(config)
	fmt.Println(req)
	if err := global.DYCLOUD_DB.Model(&configCenter.ServiceIntegration{}).Create(&req).Error; err != nil {
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
func (s *SourceCodeService) UpdateSourceCode(req *configCenter.ServiceIntegration) error {
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
	if err = global.DYCLOUD_DB.Model(&configCenter.ServiceIntegration{}).Where("id = ?", req.ID).Updates(&data).Error; err != nil {
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

	if err := global.DYCLOUD_DB.Where("id = ?", id).Delete(&configCenter.ServiceIntegration{}).Error; err != nil {
		return err
	}
	return nil
}

// DescribeSourceCode
//
//	@Description: 查看代码源详情
//	@receiver s
//	@param id
//	@return *configCenter.ServiceIntegration
//	@return error
func (s *SourceCodeService) DescribeSourceCode(id int) (*configCenter.ServiceIntegration, error) {
	fmt.Println(id)
	var data configCenter.ServiceIntegration
	if err := global.DYCLOUD_DB.Model(&configCenter.ServiceIntegration{}).Where("id = ? and type not in (0,1,2)", id).First(&data).Error; err != nil {
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
func (s *SourceCodeService) VerifySourceCode(req *configCenter.ServiceIntegration) (string, error) {
	fmt.Println(req)
	gitConf := &configCenter.GitConfig{}
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

// GetGitProjectsByRepoId
//
//	@Description: 根据仓库id查询该仓库所有项目列表
//	@receiver s
//	@param id
//	@return []*configCenter.RepoProjectRsp
//	@return error
func (s *SourceCodeService) GetGitProjectsByRepoId(id int) ([]*configCenter.RepoProjectRsp, error) {
	// 根据仓库id 获取service信息
	resp, err := s.DescribeSourceCode(id)
	if err != nil {
		return nil, err
	}
	// 定义git 配置，存储从配置中解析出来的git相关信息
	var gitConfig = &configCenter.GitConfig{}
	// 将存储在resp.Config中的JSON配置反序列化到gitConfig结构体
	err = json.Unmarshal(resp.Config, gitConfig)
	if err != nil {
		return nil, err
	}
	//   检查Git配置中是否存在Token
	if gitConfig.Token == "" {
		return nil, err
	}
	// 创建一个SCM客户端，访问SCM服务
	scmClient, err := dao.NewScmProvider(resp.Type, gitConfig.Url, gitConfig.Token)
	if err != nil {
		return nil, err
	}
	// 定义列表选项，设置分页参数
	listOptions := scm.ListOptions{
		Page: 1,
		Size: 100,
	}
	// 存储仓库列表
	repoList := []*scm.Repository{}
	// 获取第一个页面的仓库列表
	got, rsp, err := scmClient.Repositories.List(context.Background(), listOptions)
	if err != nil {
		return nil, fmt.Errorf("scmclient get repositories list error: %s", err.Error())
	}
	fmt.Println(got)
	fmt.Println(rsp)
	fmt.Println("================================================================")
	// 将第一页的仓库列表添加到repoList
	repoList = append(repoList, got...)
	// 遍历剩余的页面，从第二页到最后一页
	for i := 1; i < rsp.Page.Last; {
		// 增加页面页数
		listOptions.Page++
		// 获取下一页的仓库列表
		got, _, err := scmClient.Repositories.List(context.Background(), listOptions)
		if err != nil {
			return nil, fmt.Errorf("when get repositories list from gitlab occur error: %s", err.Error())
		}
		// 将获取到的仓库列表追加到repoList
		repoList = append(repoList, got...)
		i++ // 更新计数器
	}
	fmt.Println(repoList)
	// 创建响应结构体，存储返回数据
	newRsp := []*configCenter.RepoProjectRsp{}
	// 遍历每一个仓库，将信息转换为自定义结构
	for _, item := range repoList {
		newItem := &configCenter.RepoProjectRsp{
			RepoID:   id,                               // 仓库Id
			Path:     item.Clone,                       // 仓库的克隆路径
			FullName: item.Namespace + "/" + item.Name, // 完整名称，包含命名空间和仓库名
			Name:     item.Name,                        // 仓库名
		}
		newRsp = append(newRsp, newItem)
	}
	fmt.Println(newRsp)
	return newRsp, nil
}
