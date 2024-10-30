package ast

import (
	"DYCLOUD/global"
	"path/filepath"
)

func init() {
	global.DYCLOUD_CONFIG.AutoCode.Root, _ = filepath.Abs("../../../")
	global.DYCLOUD_CONFIG.AutoCode.Server = "server"
}
