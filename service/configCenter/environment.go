package configCenter

import (
	"DYCLOUD/global"
	"DYCLOUD/model/configCenter"
	"DYCLOUD/model/configCenter/request"
	"fmt"
)

type EnvironmentService struct{}

// GetEnvironmentList
//
//	@Description: 获取应用环境列表
//	@receiver e
//	@param req
//	@return envList
//	@return err
func (e *EnvironmentService) GetEnvironmentList(req *request.EnvRequest) (envList *[]configCenter.Environment, total int64, err error) {

	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	db := global.DYCLOUD_DB.Model(&configCenter.Environment{})

	// 创建db
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		db = db.Where("name like ?", keyword).Or("`key` like ?", keyword).Or("id = ?", req.Keyword)
	}

	var data []configCenter.Environment
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

// DescribeEnvironment
//
//	@Description: 查看环境详情
//	@receiver e
//	@param id
//	@return envList
//	@return err
func (e *EnvironmentService) DescribeEnvironment(id int) (envList *configCenter.Environment, err error) {
	fmt.Println(id)
	var data *configCenter.Environment
	if err := global.DYCLOUD_DB.Where("id = ?", id).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// CreateEnvironment
//
//	@Description: 创建环境
//	@receiver e
//	@param req
//	@return error
func (e *EnvironmentService) CreateEnvironment(req *configCenter.Environment) error {
	fmt.Println(req)
	var existingEnv configCenter.Environment
	// 检查是否存在相同的 name
	if err := global.DYCLOUD_DB.Where("name = ?", req.Name).First(&existingEnv).Error; err == nil {
		return fmt.Errorf("环境名称 '%s' 已存在，请选择一个唯一的名称", req.Name)
	}
	// 检查是否存在相同的 key
	if err := global.DYCLOUD_DB.Where("`key` = ?", req.Key).First(&existingEnv).Error; err == nil {
		return fmt.Errorf("环境标识 '%s' 已存在，请选择一个唯一的标识", req.Key)
	}
	// 如果不存在重复 key，继续创建
	if err := global.DYCLOUD_DB.Create(&req).Error; err != nil {
		return err
	}
	return nil
}

// UpdateEnvironment
//
//	@Description: 更新环境
//	@receiver e
//	@param req
//	@return data
//	@return err
func (e *EnvironmentService) UpdateEnvironment(req *configCenter.Environment) (data *configCenter.Environment, err error) {
	fmt.Println(req)
	if err = global.DYCLOUD_DB.Model(&configCenter.Environment{}).Where("id = ?", req.ID).Omit("ID").Updates(&req).Error; err != nil {
		return nil, err
	}
	return req, nil
}

// DeleteEnvironment
//
//	@Description: 删除环境
//	@receiver e
//	@param id
//	@return error
func (e *EnvironmentService) DeleteEnvironment(id int) error {
	fmt.Println(id)
	if err := global.DYCLOUD_DB.Model(&configCenter.Environment{}).Where("id = ?", id).Delete(&configCenter.Environment{}).Error; err != nil {
		return err
	}
	return nil
}

// DeleteEnvironmentByIds
//
//	@Description: 批量删除环境
//	@receiver e
//	@param ids
//	@return error
func (e *EnvironmentService) DeleteEnvironmentByIds(ids *request.DeleteEnvByIds) error {
	fmt.Println(ids)
	if err := global.DYCLOUD_DB.Model(&configCenter.Environment{}).Where("id in ?", ids.Ids).Delete(&configCenter.Environment{}).Error; err != nil {
		return err
	}
	return nil
}
