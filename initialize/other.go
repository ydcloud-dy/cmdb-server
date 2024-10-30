package initialize

import (
	"bufio"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"os"
	"strings"

	"DYCLOUD/global"
	"DYCLOUD/utils"
)

func OtherInit() {
	dr, err := utils.ParseDuration(global.DYCLOUD_CONFIG.JWT.ExpiresTime)
	if err != nil {
		panic(err)
	}
	_, err = utils.ParseDuration(global.DYCLOUD_CONFIG.JWT.BufferTime)
	if err != nil {
		panic(err)
	}

	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(dr),
	)
	file, err := os.Open("go.mod")
	if err == nil && global.DYCLOUD_CONFIG.AutoCode.Module == "" {
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		global.DYCLOUD_CONFIG.AutoCode.Module = strings.TrimPrefix(scanner.Text(), "module ")
	}
}
