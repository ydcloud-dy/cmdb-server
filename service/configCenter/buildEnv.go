package configCenter

import (
	"DYCLOUD/global"
	"DYCLOUD/model/configCenter"
	"DYCLOUD/model/configCenter/request"
	"fmt"
)

type BuildEnvService struct{}

// GetBuildEnvList
//
//	@Description: 获取应用环境列表
//	@receiver e
//	@param req
//	@return envList
//	@return err
func (e *BuildEnvService) GetBuildEnvList(req *request.BuildEnvRequest) (envList *[]configCenter.BuildEnv, total int64, err error) {

	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	db := global.DYCLOUD_DB.Model(&configCenter.BuildEnv{})

	// 创建db
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		db = db.Where("name like ?", keyword).Or("id = ?", req.Keyword)
	}

	var data []configCenter.BuildEnv
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, nil
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

// DescribeBuildEnv
//
//	@Description: 查看环境详情
//	@receiver e
//	@param id
//	@return envList
//	@return err
func (e *BuildEnvService) DescribeBuildEnv(id int) (envList *configCenter.BuildEnv, err error) {
	var data *configCenter.BuildEnv
	if err := global.DYCLOUD_DB.Where("id = ?", id).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// CreateBuildEnv
//
//	@Description: 创建环境
//	@receiver e
//	@param req
//	@return error
func (e *BuildEnvService) CreateBuildEnv(req *configCenter.BuildEnv) error {
	var existingEnv configCenter.BuildEnv
	// 检查是否存在相同的 name
	if err := global.DYCLOUD_DB.Where("name = ?", req.Name).First(&existingEnv).Error; err == nil {
		return fmt.Errorf("环境名称 '%s' 已存在，请选择一个唯一的名称", req.Name)
	}
	// 如果不存在重复 key，继续创建
	if err := global.DYCLOUD_DB.Create(&req).Error; err != nil {
		return err
	}
	return nil
}

// UpdateBuildEnv
//
//	@Description: 更新环境
//	@receiver e
//	@param req
//	@return data
//	@return err
func (e *BuildEnvService) UpdateBuildEnv(req *configCenter.BuildEnv) (data *configCenter.BuildEnv, err error) {
	if err = global.DYCLOUD_DB.Model(&configCenter.BuildEnv{}).Where("id = ?", req.ID).Omit("ID").Updates(&req).Error; err != nil {
		return nil, err
	}
	return req, nil
}

// DeleteBuildEnv
//
//	@Description: 删除环境
//	@receiver e
//	@param id
//	@return error
func (e *BuildEnvService) DeleteBuildEnv(id int) error {
	if err := global.DYCLOUD_DB.Model(&configCenter.BuildEnv{}).Where("id = ?", id).Delete(&configCenter.BuildEnv{}).Error; err != nil {
		return err
	}
	return nil
}

// DeleteBuildEnvByIds
//
//	@Description: 批量删除环境
//	@receiver e
//	@param ids
//	@return error
func (e *BuildEnvService) DeleteBuildEnvByIds(ids *request.DeleteBuildEnvByIds) error {
	if err := global.DYCLOUD_DB.Model(&configCenter.BuildEnv{}).Where("id in ?", ids.Ids).Delete(&configCenter.BuildEnv{}).Error; err != nil {
		return err
	}
	return nil
}
