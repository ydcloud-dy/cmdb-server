package initialize

import (
	"os"

	"DYCLOUD/global"
	"DYCLOUD/model/example"
	"DYCLOUD/model/system"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	switch global.DYCLOUD_CONFIG.System.DbType {
	case "mysql":
		global.DYCLOUD_ACTIVE_DBNAME = &global.DYCLOUD_CONFIG.Mysql.Dbname
		return GormMysql()
	case "pgsql":
		global.DYCLOUD_ACTIVE_DBNAME = &global.DYCLOUD_CONFIG.Pgsql.Dbname
		return GormPgSql()
	case "oracle":
		global.DYCLOUD_ACTIVE_DBNAME = &global.DYCLOUD_CONFIG.Oracle.Dbname
		return GormOracle()
	case "mssql":
		global.DYCLOUD_ACTIVE_DBNAME = &global.DYCLOUD_CONFIG.Mssql.Dbname
		return GormMssql()
	case "sqlite":
		global.DYCLOUD_ACTIVE_DBNAME = &global.DYCLOUD_CONFIG.Sqlite.Dbname
		return GormSqlite()
	default:
		global.DYCLOUD_ACTIVE_DBNAME = &global.DYCLOUD_CONFIG.Mysql.Dbname
		return GormMysql()
	}
}

func RegisterTables() {
	db := global.DYCLOUD_DB
	err := db.AutoMigrate(

		system.SysApi{},
		system.SysIgnoreApi{},
		system.SysUser{},
		system.SysBaseMenu{},
		system.JwtBlacklist{},
		system.SysAuthority{},
		system.SysDictionary{},
		system.SysOperationRecord{},
		system.SysDictionaryDetail{},
		system.SysBaseMenuParameter{},
		system.SysBaseMenuBtn{},
		system.SysAuthorityBtn{},
		system.SysExportTemplate{},
		system.Condition{},
		system.JoinTemplate{},

		example.ExaFile{},
		example.ExaCustomer{},
		example.ExaFileChunk{},
		example.ExaFileUploadAndDownload{},
	)
	if err != nil {
		global.DYCLOUD_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}

	err = bizModel()

	if err != nil {
		global.DYCLOUD_LOG.Error("register biz_table failed", zap.Error(err))
		os.Exit(0)
	}
	global.DYCLOUD_LOG.Info("register table success")
}
