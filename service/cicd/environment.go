package cicd

import (
	"DYCLOUD/global"
	"DYCLOUD/model/cicd"
	"DYCLOUD/model/cicd/request"
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
func (e *EnvironmentService) GetEnvironmentList(req *request.EnvRequest) (envList *[]cicd.Environment, total int64, err error) {

	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	db := global.DYCLOUD_DB.Model(&cicd.Environment{})

	// 创建db
	if req.Keyword != "" {
		db = db.Where("name = ?", req.Keyword).Or("id = ?", req.Keyword)
	}

	var data []cicd.Environment
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
func (e *EnvironmentService) DescribeEnvironment(id int) (envList *cicd.Environment, err error) {
	fmt.Println(id)
	var data *cicd.Environment
	if err := global.DYCLOUD_DB.Where("id = ?", id).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
func (e *EnvironmentService) CreateEnvironment(req *cicd.Environment) error {
	fmt.Println(req)

	if err := global.DYCLOUD_DB.Create(&req).Error; err != nil {
		return err
	}
	return nil
}
func (e *EnvironmentService) UpdateEnvironment(req *cicd.Environment) (data *cicd.Environment, err error) {
	fmt.Println(req)
	if err = global.DYCLOUD_DB.Model(&cicd.Environment{}).Where("id = ?", req.ID).Omit("ID").Updates(&req).Error; err != nil {
		return nil, err
	}
	return req, nil
}
func (e *EnvironmentService) DeleteEnvironment(id int) error {
	fmt.Println(id)
	if err := global.DYCLOUD_DB.Model(&cicd.Environment{}).Where("id = ?", id).Delete(&cicd.Environment{}).Error; err != nil {
		return err
	}
	return nil
}

func (e *EnvironmentService) DeleteEnvironmentByIds(ids *request.DeleteEnvByIds) error {
	fmt.Println(ids)
	if err := global.DYCLOUD_DB.Model(&cicd.Environment{}).Where("id in ?", ids.Ids).Delete(&cicd.Environment{}).Error; err != nil {
		return err
	}
	return nil
}
