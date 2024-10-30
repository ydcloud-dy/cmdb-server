package initialize

import (
	_ "DYCLOUD/source/example"
	_ "DYCLOUD/source/system"
)

func init() {
	// do nothing,only import source package so that inits can be registered
}
