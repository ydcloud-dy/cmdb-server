package global

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/qmgo"

	"DYCLOUD/utils/timer"
	"github.com/songzhibin97/gkit/cache/local_cache"

	"golang.org/x/sync/singleflight"

	"go.uber.org/zap"

	"DYCLOUD/config"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	DYCLOUD_DB        *gorm.DB
	DYCLOUD_DBList    map[string]*gorm.DB
	DYCLOUD_REDIS     redis.UniversalClient
	DYCLOUD_REDISList map[string]redis.UniversalClient
	DYCLOUD_MONGO     *qmgo.QmgoClient
	DYCLOUD_CONFIG    config.Server
	DYCLOUD_VP        *viper.Viper
	// DYCLOUD_LOG    *oplogging.Logger
	DYCLOUD_LOG                 *zap.Logger
	DYCLOUD_Timer               timer.Timer = timer.NewTimerTask()
	DYCLOUD_Concurrency_Control             = &singleflight.Group{}
	DYCLOUD_ROUTERS             gin.RoutesInfo
	DYCLOUD_ACTIVE_DBNAME       *string
	BlackCache                  local_cache.Cache
	lock                        sync.RWMutex
)

// GetGlobalDBByDBName 通过名称获取db list中的db
func GetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return DYCLOUD_DBList[dbname]
}

// MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := DYCLOUD_DBList[dbname]
	if !ok || db == nil {
		panic("db no init")
	}
	return db
}

func GetRedis(name string) redis.UniversalClient {
	redis, ok := DYCLOUD_REDISList[name]
	if !ok || redis == nil {
		panic(fmt.Sprintf("redis `%s` no init", name))
	}
	return redis
}
