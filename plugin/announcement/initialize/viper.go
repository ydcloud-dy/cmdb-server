package initialize

import (
	"DYCLOUD/global"
	"DYCLOUD/plugin/announcement/plugin"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Viper() {
	err := global.DYCLOUD_VP.UnmarshalKey("announcement", &plugin.Config)
	if err != nil {
		err = errors.Wrap(err, "初始化配置文件失败!")
		zap.L().Error(fmt.Sprintf("%+v", err))
	}
}
