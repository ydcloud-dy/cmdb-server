package service

import (
	"DYCLOUD/global"
	model "DYCLOUD/model/cloudCmdb"
	cloudcmdbreq "DYCLOUD/model/cloudCmdb/cloudcmdb"
	"DYCLOUD/model/common/request"
	"DYCLOUD/utils/cloudCmdb/aliyun"
	"DYCLOUD/utils/cloudCmdb/huawei"
	"DYCLOUD/utils/cloudCmdb/tencent"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CloudRDSService struct{}

//@function: List
//@description: RDS列表
//@param: machine model.RDS, info cloudcmdbreq.PageInfo, order string, desc bool
//@return: list interface{}, total int64, err error

func (r *CloudRDSService) List(slb model.RDS, info cloudcmdbreq.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.DYCLOUD_DB.Model(model.RDS{})
	var slbList []model.RDS

	if info.Keyword != "" && info.Field != "" {
		db = db.Where(info.Field+" LIKE ?", "%"+info.Keyword+"%")
	}

	if slb.InstanceId != "" {
		db = db.Where("instance_id LIKE ?", "%"+slb.InstanceId+"%")
	}

	if slb.Name != "" {
		db = db.Where("name = ?", slb.Name)
	}

	if slb.Region != "" {
		db = db.Where("region = ?", slb.Region)
	}

	err = db.Count(&total).Error

	if err != nil {
		return slbList, total, err
	} else {
		db = db.Limit(limit).Offset(offset)
		if order != "" {
			var OrderStr string
			// 设置有效排序key 防止sql注入
			orderMap := make(map[string]bool, 5)
			orderMap["id"] = true
			orderMap["instance_id"] = true
			orderMap["name"] = true
			orderMap["status"] = true
			if orderMap[order] {
				if desc {
					OrderStr = order + " desc"
				} else {
					OrderStr = order
				}
			} else { // didn't matched any order key in `orderMap`
				err = fmt.Errorf("非法的排序字段: %v", order)
				return slbList, total, err
			}

			err = db.Order(OrderStr).Find(&slbList).Error
		} else {
			err = db.Order("id").Find(&slbList).Error
		}
	}

	return slbList, total, err
}

//@function: UpdateRDS
//@description: 更新RDS信息
//@param: list []model.RDS
//@return:

func (r *CloudRDSService) UpdateRDS(list []model.RDS) {
	db := global.DYCLOUD_DB.Model(model.RDS{})

	for _, machine := range list {
		// 开始事务
		tx := db.Begin()

		// 更新所有存在的记录，忽略不存在的记录
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{
				"name",
				"instance_id",
				"private_addr",
				"public_addr",
				"region",
				"region_name",
				"status",
				"creation_time",
				"cloud_platform_id",
			}),
		}).Create(&machine).Error; err != nil {
			global.DYCLOUD_LOG.Error("rds  messages update fail!", zap.Error(err))
			tx.Rollback()
		}

		// 插入不存在的记录
		if err := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&machine).Error; err != nil {
			global.DYCLOUD_LOG.Error("rds messages insert fail!", zap.Error(err))
			tx.Rollback()
		}

		// 提交事务
		tx.Commit()
	}
}

//@function: AliyunSyncRDS
//@description: 阿里云同步RDS
//@param: cloud model.CloudPlatform
//@return: err error

func (r *CloudRDSService) AliyunSyncLoadRDS(cloud model.CloudPlatform) (err error) {
	var regions []model.CloudRegions
	if err = global.DYCLOUD_DB.Where("cloud_platform_id = ?", cloud.ID).Find(&regions).Error; err != nil {
		return err
	}

	for _, region := range regions {
		go func(region model.CloudRegions) {
			defer func() {
				if err := recover(); err != nil {
					global.DYCLOUD_LOG.Error(fmt.Sprintf("aliyun ecs list get fail: %s", err))
				}
			}()

			rds := aliyun.NewRDS()
			list, err := rds.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.DYCLOUD_LOG.Error("aliyun LoadBalancer list get fail: ", zap.Error(err))
				return
			}

			if len(list) > 0 {
				r.UpdateRDS(list)
			}

		}(region)
	}

	return err
}

//@function: TencentSyncRDS
//@description: 腾讯云同步RDS
//@param: cloud model.CloudPlatform
//@return: err error

func (r *CloudRDSService) TencentSyncRDS(cloud model.CloudPlatform) (err error) {
	var regions []model.CloudRegions
	if err = global.DYCLOUD_DB.Where("cloud_platform_id = ?", cloud.ID).Find(&regions).Error; err != nil {
		return err
	}

	for _, region := range regions {
		go func(region model.CloudRegions) {
			defer func() {
				if err := recover(); err != nil {
					global.DYCLOUD_LOG.Error(fmt.Sprintf("Tencent slb list get fail: %s", err))
				}
			}()

			ecs := tencent.NewRDS()
			list, err := ecs.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.DYCLOUD_LOG.Error("Tencent LoadBalancer list get fail: ", zap.Error(err))
				return
			}

			if len(list) > 0 {
				r.UpdateRDS(list)
			}

		}(region)
	}
	return err
}

//@function: HuaweiSyncRDS
//@description: 华为云同步RDS
//@param: cloud model.CloudPlatform
//@return: err error

func (r *CloudRDSService) HuaweiSyncRDS(cloud model.CloudPlatform) (err error) {
	var regions []model.CloudRegions
	if err = global.DYCLOUD_DB.Where("cloud_platform_id = ?", cloud.ID).Find(&regions).Error; err != nil {
		return err
	}

	for _, region := range regions {
		go func(region model.CloudRegions) {
			defer func() {
				if err := recover(); err != nil {
					global.DYCLOUD_LOG.Error(fmt.Sprintf("huawei rds list get fail: %s", err))
				}
			}()

			rds := huawei.NewRDS()
			list, err := rds.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.DYCLOUD_LOG.Error("huawei rds list get fail: ", zap.Error(err))
				return
			}

			if len(list) > 0 {
				r.UpdateRDS(list)
			}

		}(region)
	}
	return err
}

//@function: SyncRDS
//@description: 同步各个厂商的RDS
//@param: id int
//@return: err error

func (r *CloudRDSService) SyncRDS(id int) (err error) {
	db := global.DYCLOUD_DB.Model(model.CloudPlatform{})
	var cloud model.CloudPlatform
	if err := db.Where("id = ?", id).First(&cloud).Error; err != nil {
		return err
	}

	if cloud.Platform == "aliyun" {
		if err = r.AliyunSyncLoadRDS(cloud); err != nil {
			return err
		}
	}

	if cloud.Platform == "tencent" {
		if err = r.TencentSyncRDS(cloud); err != nil {
			return err
		}
	}

	if cloud.Platform == "huawei" {
		if err = r.HuaweiSyncRDS(cloud); err != nil {
			return err
		}
	}

	return err
}

//@function: RDSTree
//@description: rds目录树
//@param: cloud model.CloudPlatform, info request.PageInfo, order string, desc bool
//@return: list interface{}, err error

func (r *CloudRDSService) RDSTree(cloud model.CloudPlatform, info request.PageInfo, order string, desc bool) (list interface{}, err error) {
	info.PageSize, info.Page = 1000, 1
	var platformTree []model.PlatformTree
	var platform CloudPlatformService
	platformList, _, err := platform.List(cloud, info, order, desc)
	if err != nil {
		return nil, err
	}

	for _, pt := range platformList {
		var rdslist []model.RDS
		var regions []model.Regions
		if err := global.DYCLOUD_DB.Table("cloud_rds").Select("DISTINCT region, region_name").Where("cloud_platform_id = ?", pt.ID).Find(&rdslist).Error; err != nil {
			global.DYCLOUD_LOG.Error("ecs machine DISTINCT fail!", zap.Error(err))
			return nil, err
		}

		if len(rdslist) > 0 {
			for _, vmRegion := range rdslist {
				regions = append(regions, model.Regions{
					ID:         vmRegion.Region,
					Name:       vmRegion.RegionName,
					RegionId:   vmRegion.Region,
					RegionName: vmRegion.RegionName,
				})
			}
		}

		platformTree = append(platformTree, model.PlatformTree{ID: pt.ID, Name: pt.Name, Region: regions})
	}

	return platformTree, err
}
