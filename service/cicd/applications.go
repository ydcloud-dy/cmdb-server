package cicd

import (
	"DYCLOUD/global"
	"DYCLOUD/model/cicd"
	request "DYCLOUD/model/cicd/request"
	configCenter2 "DYCLOUD/model/configCenter"
	"DYCLOUD/service/configCenter"
	"DYCLOUD/service/configCenter/dao"
	"encoding/json"
	"fmt"
	"github.com/drone/go-scm/scm"
	"golang.org/x/net/context"
	"gorm.io/gorm/utils"
	"strings"
)

type ApplicationsService struct{}

// GetApplicationsList
//
//	@Description: 获取应用环境列表
//	@receiver e
//	@param req
//	@return envList
//	@return err
func (e *ApplicationsService) GetApplicationsList(req *request.ApplicationRequest) (envList *[]cicd.Applications, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	db := global.DYCLOUD_DB.Model(&cicd.Applications{})

	// 创建db
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		db = db.Where("name like ?", keyword).Or("id = ?", req.Keyword)
	}
	if !req.StartCreatedAt.IsZero() && !req.EndCreatedAt.IsZero() {
		db = db.Where("created_at BETWEEN ? AND ?", req.StartCreatedAt, req.EndCreatedAt)
		db = db.Where("name = ?", req.Keyword)
	}
	var data []cicd.Applications
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&data).Error
	if err != nil {
		return nil, 0, nil
	}
	return &data, total, nil
}

// DescribeApplications
//
//	@Description: 查看应用详情
//	@receiver e
//	@param id
//	@return *cicd.Applications
//	@return error
func (e *ApplicationsService) DescribeApplications(id int) (*cicd.Applications, error) {
	var data *cicd.Applications
	if err := global.DYCLOUD_DB.Where("id = ?", id).First(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// CreateApplications
//
//	@Description: 创建应用
//	@receiver e
//	@param req
//	@return error
func (e *ApplicationsService) CreateApplications(req *cicd.Applications) error {
	fmt.Println(req)
	if err := global.DYCLOUD_DB.Create(&req).Error; err != nil {
		return err
	}

	return nil
}

// SyncBranches
//
//	@Description: 同步应用的分支信息
//	@receiver e
//	@param id
//	@return error
func (e *ApplicationsService) SyncBranches(id int) error {
	// 根据传入的应用id查询到应用详细信息
	app, err := e.DescribeApplications(id)
	if err != nil {
		return err
	}
	// 构建代码源，然后通过上面获取的应用信息中的repo_id 查询到代码源详情
	service := configCenter.SourceCodeService{}
	result, err := service.DescribeSourceCode(app.RepoId)
	if err != nil {
		return err
	}
	gitConfig := configCenter2.GitConfig{}
	json.Unmarshal(result.Config, &gitConfig)
	fmt.Println(result)
	fmt.Println(gitConfig)
	// 创建 SCM 客户端，传入 SCM 类型（如 GitLab）、路径和访问 Token
	client, err := dao.NewScmProvider(result.Type, app.Path, gitConfig.Token)
	// 用于存储从 SCM 获取的分支列表
	branchList := []*scm.Reference{}
	// 定义分页选项，初始页码为 1，页面大小为 100
	listOptions := scm.ListOptions{
		Page: 1,
		Size: 100,
	}
	// 获取 SCM 中指定应用的分支列表，返回第一页结果和分页信息
	got, res, err := client.Git.ListBranches(context.Background(), app.FullName, listOptions)
	if err != nil {
		return err
	}
	// 将获取到的分支添加到 branchList 中
	branchList = append(branchList, got...)
	// 循环处理分页数据，继续获取剩余的分支列表
	for i := 1; i < res.Page.Last; {
		// 移动到下一页
		listOptions.Page++
		// 获取下一页的分支列表
		got, _, err := client.Git.ListBranches(context.Background(), app.FullName, listOptions)
		if err != nil {
			return err
		}
		// 将获取到的分支添加到 branchList 中
		branchList = append(branchList, got...)
		// 增加页码计数
		i++
	}
	// 遍历所有获取到的分支
	for _, branch := range branchList {
		// 如果分支名称以 "release_" 开头，跳过该分支
		if strings.HasPrefix(branch.Name, "release_") {
			continue
		}
		originBranch, err := e.GetAppBranchByName(id, branch.Name)
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				err = nil
			} else {
				return fmt.Errorf("when get app branch occur error: %s", err.Error())
			}
		}
		if originBranch == nil {
			appBranch := &cicd.AppBranch{
				BranchName: branch.Name,
				Path:       app.Path,
				AppID:      id,
			}
			if _, err := e.CreateAppBranchIfNotExist(appBranch); err != nil {
				return err
			}

		} else {
			originBranch.Path = app.Path
			if err := e.UpdateAppBranch(originBranch); err != nil {
				return err
			}
		}

	}
	branchListInDB, err := e.GetAppBranches(id)
	if err != nil {
		return err
	}
	branchNameList := []string{}
	for _, branch := range branchList {
		branchNameList = append(branchNameList, branch.Name)
	}
	for _, branchDBItem := range branchListInDB {
		if !utils.Contains(branchNameList, branchDBItem.BranchName) {
			e.SoftDeleteAppBranch(branchDBItem)
		}
	}

	return nil
}
func (e *ApplicationsService) SoftDeleteAppBranch(branch *cicd.AppBranch) error {
	err := global.DYCLOUD_DB.Model(&cicd.AppBranch{}).Where("id = ?", branch.ID).Delete(&cicd.AppBranch{}).Error

	return err
}
func (e *ApplicationsService) GetAppBranches(appID int) ([]*cicd.AppBranch, error) {
	branches := []*cicd.AppBranch{}
	query := global.DYCLOUD_DB.Model(&cicd.AppBranch{})
	if appID != 0 {
		query = query.Where("app_id = ?", appID)
	}
	err := query.Find(&branches).Error
	return branches, err
}
func (e *ApplicationsService) GetBranchesList(req request.ApplicationRequest) (envList *[]cicd.AppBranch, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	db := global.DYCLOUD_DB.Model(&cicd.AppBranch{}).Where("app_id=?", req.AppId)

	var data []cicd.AppBranch
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&data).Error
	if err != nil {
		return nil, 0, nil
	}
	return &data, total, nil
}
func (e *ApplicationsService) GetAppBranchByName(appID int, branchName string) (*cicd.AppBranch, error) {
	branch := cicd.AppBranch{}
	if err := global.DYCLOUD_DB.Where("app_id=?", appID).Where("branch_name=?", branchName).First(&branch).Error; err != nil {
		return nil, err
	}
	return &branch, nil
}
func (e *ApplicationsService) CreateAppBranchIfNotExist(branch *cicd.AppBranch) (int, error) {
	result := global.DYCLOUD_DB.Where("branch_name = ? and app_id = ?", branch.BranchName, branch.AppID).FirstOrCreate(branch)
	if result.Error != nil {
		return 0, result.Error
	}
	// 检查是否是创建的新记录，或者已存在
	if result.RowsAffected == 0 {
		return int(branch.ID), fmt.Errorf("branch_name: %v already exists in app branch table", branch.BranchName)
	}
	return int(branch.ID), nil

}
func (e *ApplicationsService) UpdateAppBranch(branch *cicd.AppBranch) error {
	err := global.DYCLOUD_DB.Model(&cicd.AppBranch{}).Where("id = ?", branch.ID).Updates(branch).Error
	return err
}

// UpdateApplications
//
//	@Description: 更新应用
//	@receiver e
//	@param req
//	@return *cicd.Applications
//	@return error
func (e *ApplicationsService) UpdateApplications(req *cicd.Applications) (*cicd.Applications, error) {
	fmt.Println(req)
	data, err := e.DescribeApplications(int(req.ID))
	if err != nil {
		return nil, err
	}
	data = req
	if err = global.DYCLOUD_DB.Model(&cicd.Applications{}).Where("id = ?", req.ID).Omit("ID").Updates(&req).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// DeleteApplications
//
//	@Description: 删除应用
//	@receiver e
//	@param id
//	@return error
func (e *ApplicationsService) DeleteApplications(id int) error {
	fmt.Println(id)
	if err := global.DYCLOUD_DB.Model(&cicd.Applications{}).Where("id = ?", id).Delete(&cicd.Applications{}).Error; err != nil {
		return err
	}
	return nil
}

// DeleteApplicationsByIds
//
//	@Description: 批量删除应用
//	@receiver e
//	@param ids
//	@return error
func (e *ApplicationsService) DeleteApplicationsByIds(ids *request.DeleteApplicationByIds) error {
	fmt.Println(ids)
	if err := global.DYCLOUD_DB.Model(&cicd.Applications{}).Where("id in ?", ids.Ids).Delete(&cicd.Applications{}).Error; err != nil {
		return err
	}
	return nil
}
