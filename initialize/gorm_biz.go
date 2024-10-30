package initialize

import (
	"DYCLOUD/global"
	"DYCLOUD/model/cmdb"
)

func bizModel() error {
	db := global.DYCLOUD_DB
	err := db.AutoMigrate(cmdb.CmdbProjects{})
	if err != nil {
		return err
	}
	return nil
}
