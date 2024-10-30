package system

import (
	"DYCLOUD/global"
	"DYCLOUD/model/system/response"
)

type AutoCodeService struct{}

type Database interface {
	GetDB(businessDB string) (data []response.Db, err error)
	GetTables(businessDB string, dbName string) (data []response.Table, err error)
	GetColumn(businessDB string, tableName string, dbName string) (data []response.Column, err error)
}

func (autoCodeService *AutoCodeService) Database(businessDB string) Database {

	if businessDB == "" {
		switch global.DYCLOUD_CONFIG.System.DbType {
		case "mysql":
			return AutoCodeMysql
		case "pgsql":
			return AutoCodePgsql
		case "mssql":
			return AutoCodeMssql
		case "oracle":
			return AutoCodeOracle
		case "sqlite":
			return AutoCodeSqlite
		default:
			return AutoCodeMysql
		}
	} else {
		for _, info := range global.DYCLOUD_CONFIG.DBList {
			if info.AliasName == businessDB {
				switch info.Type {
				case "mysql":
					return AutoCodeMysql
				case "mssql":
					return AutoCodeMssql
				case "pgsql":
					return AutoCodePgsql
				case "oracle":
					return AutoCodeOracle
				case "sqlite":
					return AutoCodeSqlite
				default:
					return AutoCodeMysql
				}
			}
		}
		return AutoCodeMysql
	}

}
