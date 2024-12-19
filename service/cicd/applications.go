package cicd

import (
	"DYCLOUD/global"
	"DYCLOUD/model/cicd"
	request "DYCLOUD/model/cicd/request"
	common "DYCLOUD/model/common/request"
	"DYCLOUD/model/kubernetes/pods"
	"DYCLOUD/service/kubernetes/workload/pod"
	"fmt"
	corev1 "k8s.io/api/core/v1"
)

type ApplicationsService struct{}

// GetApplicationsList
//
//	@Description: 获取应用环境列表
//	@receiver e
//	@param req
//	@return envList
//	@return err
func (e *ApplicationsService) GetApplicationsList(req *request.ApplicationRequest) (appList *[]cicd.App, total int64, err error) {
	// 分页参数计算
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	// 初始化数据库查询
	db := global.DYCLOUD_DB.Model(&cicd.App{})

	// 条件查询：关键字搜索
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		db = db.Where("app_name LIKE ? OR app_code LIKE ? OR id = ?", keyword, keyword, req.Keyword)
	}

	// 条件查询：创建时间范围
	if !req.StartCreatedAt.IsZero() && !req.EndCreatedAt.IsZero() {
		db = db.Where("created_at BETWEEN ? AND ?", req.StartCreatedAt, req.EndCreatedAt)
	}

	// 统计总数
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("获取应用列表总数失败: %w", err)
	}

	// 分页查询
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	// 查询应用列表
	var data []cicd.App
	err = db.Find(&data).Error
	if err != nil {
		return nil, 0, fmt.Errorf("获取应用列表失败: %w", err)
	}

	// 手动加载关联数据
	for i := range data {
		// 加载环境数据
		var envs []cicd.Env
		err = global.DYCLOUD_DB.Where("app_id = ?", data[i].ID).Find(&envs).Error
		if err != nil {
			return nil, 0, fmt.Errorf("加载环境数据失败: %w", err)
		}
		data[i].Envs = envs

		// 加载开发者数据
		var developers []cicd.Developer
		err = global.DYCLOUD_DB.Where("app_id = ? AND role_type = ?", data[i].ID, "develop").Find(&developers).Error
		if err != nil {
			return nil, 0, fmt.Errorf("加载开发者数据失败: %w", err)
		}
		for j := range developers {
			developers[j].Option = &cicd.Option{
				Avatar:   developers[j].Avatar,
				Nickname: developers[j].Nickname,
				Username: developers[j].Username,
			}
		}
		data[i].Develop = developers

		// 加载拥有者数据
		var owners []cicd.Developer
		err = global.DYCLOUD_DB.Where("app_id = ? AND role_type = ?", data[i].ID, "owner").Find(&owners).Error
		if err != nil {
			return nil, 0, fmt.Errorf("加载拥有者数据失败: %w", err)
		}
		for j := range owners {
			owners[j].Option = &cicd.Option{
				Avatar:   owners[j].Avatar,
				Nickname: owners[j].Nickname,
				Username: owners[j].Username,
			}
		}
		data[i].Owner = owners
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
func (e *ApplicationsService) DescribeApplications(id int) (*cicd.App, error) {
	// 查询 App 主数据
	var app cicd.App
	if err := global.DYCLOUD_DB.Where("id = ?", id).First(&app).Error; err != nil {
		return nil, fmt.Errorf("查询应用信息失败: %w", err)
	}

	// 查询关联的 Envs 数据
	var envs []cicd.Env
	if err := global.DYCLOUD_DB.Where("app_id = ?", id).Find(&envs).Error; err != nil {
		return nil, fmt.Errorf("查询环境信息失败: %w", err)
	}
	app.Envs = envs

	// 查询关联的 Develop 数据
	var developers []cicd.Developer
	if err := global.DYCLOUD_DB.Where("app_id = ? AND role_type = ?", id, "develop").Find(&developers).Error; err != nil {
		return nil, fmt.Errorf("查询开发者信息失败: %w", err)
	}

	// 填充 Develop 的 Option 数据
	for i := range developers {
		developers[i].Option = &cicd.Option{
			Avatar:   developers[i].Avatar,
			Nickname: developers[i].Nickname,
			Username: developers[i].Username,
		}
	}
	app.Develop = developers

	// 查询关联的 Owner 数据
	var owners []cicd.Developer
	if err := global.DYCLOUD_DB.Where("app_id = ? AND role_type = ?", id, "owner").Find(&owners).Error; err != nil {
		return nil, fmt.Errorf("查询拥有者信息失败: %w", err)
	}

	// 填充 Owner 的 Option 数据
	for i := range owners {
		owners[i].Option = &cicd.Option{
			Avatar:   owners[i].Avatar,
			Nickname: owners[i].Nickname,
			Username: owners[i].Username,
		}
	}
	app.Owner = owners

	return &app, nil
}
func (e *ApplicationsService) DescribeApplicationsByName(name string) (*cicd.App, error) {
	// 查询 App 主数据
	var app cicd.App
	if err := global.DYCLOUD_DB.Where("app_code = ?", name).First(&app).Error; err != nil {
		return nil, fmt.Errorf("查询应用信息失败: %w", err)
	}

	// 查询关联的 Envs 数据
	var envs []cicd.Env
	if err := global.DYCLOUD_DB.Where("app_id = ?", app.ID).Find(&envs).Error; err != nil {
		return nil, fmt.Errorf("查询环境信息失败: %w", err)
	}
	app.Envs = envs

	// 查询关联的 Develop 数据
	var developers []cicd.Developer
	if err := global.DYCLOUD_DB.Where("app_id = ? AND role_type = ?", app.ID, "develop").Find(&developers).Error; err != nil {
		return nil, fmt.Errorf("查询开发者信息失败: %w", err)
	}

	// 填充 Develop 的 Option 数据
	for i := range developers {
		developers[i].Option = &cicd.Option{
			Avatar:   developers[i].Avatar,
			Nickname: developers[i].Nickname,
			Username: developers[i].Username,
		}
	}
	app.Develop = developers

	// 查询关联的 Owner 数据
	var owners []cicd.Developer
	if err := global.DYCLOUD_DB.Where("app_id = ? AND role_type = ?", app.ID, "owner").Find(&owners).Error; err != nil {
		return nil, fmt.Errorf("查询拥有者信息失败: %w", err)
	}

	// 填充 Owner 的 Option 数据
	for i := range owners {
		owners[i].Option = &cicd.Option{
			Avatar:   owners[i].Avatar,
			Nickname: owners[i].Nickname,
			Username: owners[i].Username,
		}
	}
	app.Owner = owners

	return &app, nil
}

// CreateApplications
//
//	@Description: 创建应用
//	@receiver e
//	@param req
//	@return error
func (e *ApplicationsService) CreateApplications(req *cicd.AppRequestBody) error {
	global.DYCLOUD_DB = global.DYCLOUD_DB.Debug()

	// 开启事务
	tx := global.DYCLOUD_DB.Begin()

	// 设置关联数据
	for i := range req.App.Develop {
		req.App.Develop[i].ID = 0               // 确保 ID 被清零
		req.App.Develop[i].RoleType = "develop" // 设置开发者角色
		if req.App.Develop[i].Option != nil {   // 提取 Option 数据
			req.App.Develop[i].Avatar = req.App.Develop[i].Option.Avatar
			req.App.Develop[i].Nickname = req.App.Develop[i].Option.Nickname
			req.App.Develop[i].Username = req.App.Develop[i].Option.Username
		}
	}
	for i := range req.App.Owner {
		req.App.Owner[i].ID = 0             // 确保 ID 被清零
		req.App.Owner[i].RoleType = "owner" // 设置拥有者角色
		if req.App.Owner[i].Option != nil { // 提取 Option 数据
			req.App.Owner[i].Avatar = req.App.Owner[i].Option.Avatar
			req.App.Owner[i].Nickname = req.App.Owner[i].Option.Nickname
			req.App.Owner[i].Username = req.App.Owner[i].Option.Username
		}
	}

	// 保存 App 数据
	if err := tx.Create(&req.App).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("保存应用信息失败: %w", err)
	}

	// 确保 App ID 已生成
	if req.App.ID == 0 {
		tx.Rollback()
		return fmt.Errorf("应用 ID 未生成")
	}

	// 设置环境数据的外键关联
	for i := range req.Envs {
		req.Envs[i].ID = 0 // 确保 ID 被清零
		req.Envs[i].AppID = req.App.ID
	}

	// 保存环境数据
	for _, env := range req.Envs {
		if err := tx.Create(&env).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("保存环境信息失败: %w", err)
		}
	}

	// 保存开发者数据
	for i := range req.App.Develop {
		req.App.Develop[i].AppID = req.App.ID // 设置外键关联
		req.App.Develop[i].ID = 0             // 确保 ID 被清零
		if err := tx.Create(&req.App.Develop[i]).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("保存开发者信息失败: %w", err)
		}
	}

	// 保存拥有者数据
	for i := range req.App.Owner {
		req.App.Owner[i].AppID = req.App.ID // 设置外键关联
		req.App.Owner[i].ID = 0             // 确保 ID 被清零
		if err := tx.Create(&req.App.Owner[i]).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("保存拥有者信息失败: %w", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("事务提交失败: %w", err)
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
	//app, err := e.DescribeApplications(id)
	//if err != nil {
	//	return err
	//}
	//// 构建代码源，然后通过上面获取的应用信息中的repo_id 查询到代码源详情
	//service := configCenter.SourceCodeService{}
	//result, err := service.DescribeSourceCode(app.RepoId)
	//if err != nil {
	//	return err
	//}
	//
	//gitConfig := configCenter2.GitConfig{}
	//json.Unmarshal(result.Config, &gitConfig)
	//fmt.Println(result)
	//fmt.Println(gitConfig)
	//// 创建 SCM 客户端，传入 SCM 类型（如 GitLab）、路径和访问 Token
	//client, err := dao.NewScmProvider(result.Type, app.Path, gitConfig.Token)
	//// 用于存储从 SCM 获取的分支列表
	//branchList := []*scm.Reference{}
	//// 定义分页选项，初始页码为 1，页面大小为 100
	//listOptions := scm.ListOptions{
	//	Page: 1,
	//	Size: 100,
	//}
	//// 获取 SCM 中指定应用的分支列表，返回第一页结果和分页信息
	//got, res, err := client.Git.ListBranches(context.Background(), app.FullName, listOptions)
	//if err != nil {
	//	return err
	//}
	//// 将获取到的分支添加到 branchList 中
	//branchList = append(branchList, got...)
	//// 循环处理分页数据，继续获取剩余的分支列表
	//for i := 1; i < res.Page.Last; {
	//	// 移动到下一页
	//	listOptions.Page++
	//	// 获取下一页的分支列表
	//	got, _, err := client.Git.ListBranches(context.Background(), app.FullName, listOptions)
	//	if err != nil {
	//		return err
	//	}
	//	// 将获取到的分支添加到 branchList 中
	//	branchList = append(branchList, got...)
	//	// 增加页码计数
	//	i++
	//}
	//// 遍历所有获取到的分支
	//for _, branch := range branchList {
	//	// 如果分支名称以 "release_" 开头，跳过该分支
	//	if strings.HasPrefix(branch.Name, "release_") {
	//		continue
	//	}
	//	originBranch, err := e.GetAppBranchByName(id, branch.Name)
	//	if err != nil {
	//		if strings.Contains(err.Error(), "record not found") {
	//			err = nil
	//		} else {
	//			return fmt.Errorf("when get app branch occur error: %s", err.Error())
	//		}
	//	}
	//	if originBranch == nil {
	//		appBranch := &cicd.AppBranch{
	//			BranchName: branch.Name,
	//			Path:       app.Path,
	//			AppID:      id,
	//		}
	//		if _, err := e.CreateAppBranchIfNotExist(appBranch); err != nil {
	//			return err
	//		}
	//
	//	} else {
	//		originBranch.Path = app.Path
	//		if err := e.UpdateAppBranch(originBranch); err != nil {
	//			return err
	//		}
	//	}
	//
	//}
	//branchListInDB, err := e.GetAppBranches(id)
	//if err != nil {
	//	return err
	//}
	//branchNameList := []string{}
	//for _, branch := range branchList {
	//	branchNameList = append(branchNameList, branch.Name)
	//}
	//for _, branchDBItem := range branchListInDB {
	//	if !utils.Contains(branchNameList, branchDBItem.BranchName) {
	//		e.SoftDeleteAppBranch(branchDBItem)
	//	}
	//}

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
func (e *ApplicationsService) UpdateApplications(req *cicd.AppRequestBody) (*cicd.AppRequestBody, error) {
	global.DYCLOUD_DB = global.DYCLOUD_DB.Debug()
	// 开启事务
	tx := global.DYCLOUD_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	// 更新 App 数据
	app := req.App
	if err := tx.Model(&cicd.App{}).Where("id = ?", req.App.ID).Updates(app).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新应用信息失败: %w", err)
	}

	// 确保应用 ID 存在
	if app.ID == 0 {
		tx.Rollback()
		return nil, fmt.Errorf("应用 ID 未生成或无效")
	}
	// **更新 Envs 数据**
	// 获取当前环境 ID 列表
	envIDs := make([]uint, 0)
	for _, env := range req.Envs {
		if env.ID != 0 {
			envIDs = append(envIDs, env.ID)
		}
	}
	// 删除不在请求中的环境记录
	if err := tx.Where("app_id = ? AND id NOT IN ?", app.ID, envIDs).Delete(&cicd.Env{}).Error; err != nil {
		return nil, fmt.Errorf("清理环境信息失败: %w", err)
	}

	// 更新 Envs 数据
	for _, env := range req.Envs {
		env.AppID = app.ID
		if env.ID == 0 { // 新增
			if err := tx.Create(&env).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("新增环境信息失败: %w", err)
			}
		} else { // 更新
			if err := tx.Model(&cicd.Env{}).Where("id = ?", env.ID).Updates(env).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("更新环境信息失败: %w", err)
			}
		}
	}
	// **更新 Develop 数据**
	// 删除原有 Develop 数据
	if err := tx.Where("app_id = ? AND role_type = ?", app.ID, "develop").Delete(&cicd.Developer{}).Error; err != nil {
		return nil, fmt.Errorf("清理开发者信息失败: %w", err)
	}

	// 更新 Develop 数据
	for _, dev := range app.Develop {
		dev.AppID = app.ID
		dev.RoleType = "develop"
		dev.ID = 0 // 确保重新插入
		if dev.Option != nil {
			dev.Avatar = dev.Option.Avatar
			dev.Nickname = dev.Option.Nickname
			dev.Username = dev.Option.Username
		}
		if err := tx.Create(&dev).Error; err != nil {
			return nil, fmt.Errorf("新增开发者信息失败: %w", err)
		}
	}
	// **更新 Owner 数据**
	// 删除原有 Owner 数据
	if err := tx.Where("app_id = ? AND role_type = ?", app.ID, "owner").Delete(&cicd.Developer{}).Error; err != nil {
		return nil, fmt.Errorf("清理拥有者信息失败: %w", err)
	}
	// 插入新的 Owner 数据
	for _, owner := range app.Owner {
		owner.AppID = app.ID
		owner.RoleType = "owner"
		owner.ID = 0 // 确保重新插入
		if owner.Option != nil {
			owner.Avatar = owner.Option.Avatar
			owner.Nickname = owner.Option.Nickname
			owner.Username = owner.Option.Username
		}
		if err := tx.Create(&owner).Error; err != nil {
			return nil, fmt.Errorf("新增拥有者信息失败: %w", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("事务提交失败: %w", err)
	}

	return req, nil
}

// DeleteApplications
//
//	@Description: 删除应用
//	@receiver e
//	@param id
//	@return error
func (e *ApplicationsService) DeleteApplications(id int) error {
	// 开启事务
	tx := global.DYCLOUD_DB.Begin()

	// 删除关联的开发者数据 (Develop)
	if err := tx.Where("app_id = ? AND role_type = ?", id, "develop").Delete(&cicd.Developer{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除开发者信息失败: %w", err)
	}

	// 删除关联的拥有者数据 (Owner)
	if err := tx.Where("app_id = ? AND role_type = ?", id, "owner").Delete(&cicd.Developer{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除拥有者信息失败: %w", err)
	}

	// 删除关联的环境数据 (Envs)
	if err := tx.Where("app_id = ?", id).Delete(&cicd.Env{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除环境信息失败: %w", err)
	}

	// 删除主表 App 数据
	if err := tx.Where("id = ?", id).Delete(&cicd.App{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除应用信息失败: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("事务提交失败: %w", err)
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

func (e *ApplicationsService) GetApplicationDeploymentInfo(req *request.DeploymentInfoRequest) (*[]corev1.Pod, int, error) {
	fmt.Println(req)
	service := pod.K8sPodService{}
	podList, total, err := service.GetPodList(pods.PodListReq{
		ClusterId:     req.ClusterId,
		Namespace:     req.Namespace,
		LabelSelector: req.AppCode,
		PageInfo: common.PageInfo{
			Page:     1,
			PageSize: 10000,
		},
	})
	if err != nil {
		return nil, 0, err
	}
	fmt.Println(podList)
	fmt.Println(total)
	return podList, total, nil
}
