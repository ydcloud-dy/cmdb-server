package cmdb

import (
	"DYCLOUD/global"
	"DYCLOUD/model/cmdb"
    cmdbReq "DYCLOUD/model/cmdb/request"
    "gorm.io/gorm"
)

type CmdbProjectsService struct {}
// CreateCmdbProjects 创建cmdbProjects表记录
// Author [yourname](https://github.com/yourname)
func (cmdbProjectsService *CmdbProjectsService) CreateCmdbProjects(cmdbProjects *cmdb.CmdbProjects) (err error) {
	err = global.DYCLOUD_DB.Create(cmdbProjects).Error
	return err
}

// DeleteCmdbProjects 删除cmdbProjects表记录
// Author [yourname](https://github.com/yourname)
func (cmdbProjectsService *CmdbProjectsService)DeleteCmdbProjects(ID string,userID uint) (err error) {
	err = global.DYCLOUD_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&cmdb.CmdbProjects{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
              return err
        }
        if err = tx.Delete(&cmdb.CmdbProjects{},"id = ?",ID).Error; err != nil {
              return err
        }
        return nil
	})
	return err
}

// DeleteCmdbProjectsByIds 批量删除cmdbProjects表记录
// Author [yourname](https://github.com/yourname)
func (cmdbProjectsService *CmdbProjectsService)DeleteCmdbProjectsByIds(IDs []string,deleted_by uint) (err error) {
	err = global.DYCLOUD_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&cmdb.CmdbProjects{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
            return err
        }
        if err := tx.Where("id in ?", IDs).Delete(&cmdb.CmdbProjects{}).Error; err != nil {
            return err
        }
        return nil
    })
	return err
}

// UpdateCmdbProjects 更新cmdbProjects表记录
// Author [yourname](https://github.com/yourname)
func (cmdbProjectsService *CmdbProjectsService)UpdateCmdbProjects(cmdbProjects cmdb.CmdbProjects) (err error) {
	err = global.DYCLOUD_DB.Model(&cmdb.CmdbProjects{}).Where("id = ?",cmdbProjects.ID).Updates(&cmdbProjects).Error
	return err
}

// GetCmdbProjects 根据ID获取cmdbProjects表记录
// Author [yourname](https://github.com/yourname)
func (cmdbProjectsService *CmdbProjectsService)GetCmdbProjects(ID string) (cmdbProjects cmdb.CmdbProjects, err error) {
	err = global.DYCLOUD_DB.Where("id = ?", ID).First(&cmdbProjects).Error
	return
}

// GetCmdbProjectsInfoList 分页获取cmdbProjects表记录
// Author [yourname](https://github.com/yourname)
func (cmdbProjectsService *CmdbProjectsService)GetCmdbProjectsInfoList(info cmdbReq.CmdbProjectsSearch) (list []cmdb.CmdbProjects, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.DYCLOUD_DB.Model(&cmdb.CmdbProjects{})
    var cmdbProjectss []cmdb.CmdbProjects
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&cmdbProjectss).Error
	return  cmdbProjectss, total, err
}
func (cmdbProjectsService *CmdbProjectsService)GetCmdbProjectsPublic() {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
