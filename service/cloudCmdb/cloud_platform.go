package service

import (
	"DYCLOUD/global"
	model "DYCLOUD/model/cloudCmdb"
	"DYCLOUD/model/common/request"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CloudPlatformService struct{}

//@function: List
//@description: 厂商列表
//@param: cloud model.CloudPlatform, info request.PageInfo, order string, desc bool
//@return: list []model.CloudPlatform, total int64, err error

func (p *CloudPlatformService) List(cloud model.CloudPlatform, info request.PageInfo, order string, desc bool) (list []model.CloudPlatform, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.DYCLOUD_DB.Model(model.CloudPlatform{})
	var cloudList []model.CloudPlatform

	if cloud.Name != "" {
		db = db.Where("name LIKE ?", "%"+cloud.Name+"%")
	}

	err = db.Count(&total).Error

	if err != nil {
		return cloudList, total, err
	} else {
		db = db.Limit(limit).Offset(offset)
		if order != "" {
			var OrderStr string
			// 设置有效排序key 防止sql注入
			orderMap := make(map[string]bool, 5)
			orderMap["id"] = true
			orderMap["name"] = true
			if orderMap[order] {
				if desc {
					OrderStr = order + " desc"
				} else {
					OrderStr = order
				}
			} else { // didn't matched any order key in `orderMap`
				err = fmt.Errorf("非法的排序字段: %v", order)
				return cloudList, total, err
			}

			err = db.Order(OrderStr).Find(&cloudList).Error
		} else {
			err = db.Order("id").Find(&cloudList).Error
		}
	}

	return cloudList, total, err
}

//@function: GetCloudPlatformById
//@description: 获取单个厂商信息
//@param: id int
//@return: cloud model.CloudPlatform, regions []model.CloudRegions, err error

func (p *CloudPlatformService) GetCloudPlatformById(id int) (cloud model.CloudPlatform, regions []model.CloudRegions, err error) {
	if err = global.DYCLOUD_DB.Where("id = ?", id).First(&cloud).Error; err != nil {
		return cloud, regions, nil
	}

	if err = global.DYCLOUD_DB.Select("id, name").Where("cloud_platform_id = ?", id).Find(&regions).Error; err != nil {
		return cloud, regions, nil
	}

	return cloud, regions, nil
}

//@function: CreateCloudPlatform
//@description: 创建厂商信息
//@param: cloud model.CloudPlatform
//@return: err error

func (p *CloudPlatformService) CreateCloudPlatform(cloud model.CloudPlatform) (err error) {
	if !errors.Is(global.DYCLOUD_DB.Where("name = ?", cloud.Name).First(&model.CloudPlatform{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同厂商")
	}
	return global.DYCLOUD_DB.Create(&cloud).Error
}

//@function: UpdateCloudPlatform
//@description: 更新厂商信息
//@param: cloud model.CloudPlatform
//@return: err error

func (p *CloudPlatformService) UpdateCloudPlatform(cloud model.CloudPlatform) (err error) {
	return global.DYCLOUD_DB.Where("id = ?", cloud.ID).First(&model.CloudPlatform{}).Updates(&cloud).Error
}

//@function: DeleteCloudPlatform
//@description: 删除厂商信息
//@param: req request.GetById
//@return: err error

func (p *CloudPlatformService) DeleteCloudPlatform(req request.GetById) (err error) {
	var cloud model.CloudPlatform
	if err = global.DYCLOUD_DB.Where("id = ?", req.ID).First(&cloud).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err = global.DYCLOUD_DB.Delete(&cloud).Error; err != nil {
		return err
	}

	return err
}

//@function: DeleteCloudPlatformByIds
//@description: 批量删除厂商信息
//@param: ids request.IdsReq
//@return: err error

func (p *CloudPlatformService) DeleteCloudPlatformByIds(ids request.IdsReq) (err error) {
	return global.DYCLOUD_DB.Delete(&[]model.CloudPlatform{}, "id in ?", ids.Ids).Error
}
