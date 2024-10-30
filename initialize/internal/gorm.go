package internal

import (
	"DYCLOUD/config"
	"DYCLOUD/global"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var Gorm = new(_gorm)

type _gorm struct{}

// Config gorm 自定义配置
// Author [SliverHorn](https://github.com/SliverHorn)
func (g *_gorm) Config(prefix string, singular bool) *gorm.Config {
	var general config.GeneralDB
	switch global.DYCLOUD_CONFIG.System.DbType {
	case "mysql":
		general = global.DYCLOUD_CONFIG.Mysql.GeneralDB
	case "pgsql":
		general = global.DYCLOUD_CONFIG.Pgsql.GeneralDB
	case "oracle":
		general = global.DYCLOUD_CONFIG.Oracle.GeneralDB
	case "sqlite":
		general = global.DYCLOUD_CONFIG.Sqlite.GeneralDB
	case "mssql":
		general = global.DYCLOUD_CONFIG.Mssql.GeneralDB
	default:
		general = global.DYCLOUD_CONFIG.Mysql.GeneralDB
	}
	return &gorm.Config{
		Logger: logger.New(NewWriter(general, log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      general.LogLevel(),
			Colorful:      true,
		}),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
			SingularTable: singular,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}
}
