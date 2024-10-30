package service

import (
	"DYCLOUD/global"
	global2 "DYCLOUD/global/docker"
	"DYCLOUD/model/docker"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type HostService struct {
}

// CreateHost 创建主机列表记录

func (hsService *HostService) CreateHost(hs *model.Host) (err error) {

	err = global.DYCLOUD_DB.Create(hs).Error
	return err
}

// DeleteHost 删除主机列表记录

func (hsService *HostService) DeleteHost(ID string, userID uint) (err error) {
	err = global.DYCLOUD_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Host{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&model.Host{}, "id = ?", ID).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}

// DeleteHostByIds 批量删除主机列表记录
func (hsService *HostService) DeleteHostByIds(IDs []string, deleted_by uint) (err error) {
	err = global.DYCLOUD_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Host{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", IDs).Delete(&model.Host{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateHost 更新主机列表记录
func (hsService *HostService) UpdateHost(hs model.Host) (err error) {
	err = global.DYCLOUD_DB.Model(&model.Host{}).Where("id = ?", hs.ID).Updates(&hs).Error
	global2.DockerClient.Remove(hs.Name)
	return err
}

// GetHost 根据ID获取主机列表记录

func (hsService *HostService) GetHost(ID string) (hs model.Host, err error) {
	err = global.DYCLOUD_DB.Where("id = ?", ID).First(&hs).Error
	return
}

// GetHostInfoList 分页获取主机列表记录

func (hsService *HostService) GetHostInfoList(info model.HostSearch) (list []model.Host, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.DYCLOUD_DB.Model(&model.Host{})
	var hss []model.Host

	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&hss).Error
	//hss = append(hss, model.Host{DYCLOUD_MODEL: global.DYCLOUD_MODEL{
	//	CreatedAt: time.Now(),
	//}, Name: "host", Description: "宿主机"})
	return hss, total, err
}

// CheckHost 测试host主机可用性

func (hsService *HostService) CheckHost(hs model.Host) (hsRes *model.HostCheckResponse, err error) {
	hsRes = &model.HostCheckResponse{
		Success: false,
	}
	client, err := global2.DockerClient.CreateClient(hs)
	if err != nil {
		return nil, err
	}
	res, err := client.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println(res)
	hsRes.Success = true

	return hsRes, err
}
