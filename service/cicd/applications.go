package cicd

import (
	"DYCLOUD/global"
	"DYCLOUD/model/cicd"
	request "DYCLOUD/model/cicd/request"
	"fmt"
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
func (e *ApplicationsService) DescribeApplications(id int) (*cicd.Applications, error) {
	var data *cicd.Applications
	if err := global.DYCLOUD_DB.Where("id = ?", id).First(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
func (e *ApplicationsService) CreateApplications(req *cicd.Applications) error {
	fmt.Println(req)
	if err := global.DYCLOUD_DB.Create(&req).Error; err != nil {
		return err
	}

	return nil
}
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
func (e *ApplicationsService) DeleteApplications(id int) error {
	fmt.Println(id)
	if err := global.DYCLOUD_DB.Model(&cicd.Applications{}).Where("id = ?", id).Delete(&cicd.Applications{}).Error; err != nil {
		return err
	}
	return nil
}

func (e *ApplicationsService) DeleteApplicationsByIds(ids *request.DeleteApplicationByIds) error {
	fmt.Println(ids)
	if err := global.DYCLOUD_DB.Model(&cicd.Applications{}).Where("id in ?", ids.Ids).Delete(&cicd.Applications{}).Error; err != nil {
		return err
	}
	return nil
}
