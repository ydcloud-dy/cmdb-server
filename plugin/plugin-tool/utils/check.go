package utils

import (
	"fmt"

	"DYCLOUD/global"
	"DYCLOUD/model/system"
)

func RegisterApis(apis ...system.SysApi) {
	var count int64
	var apiPaths []string
	for i := range apis {
		apiPaths = append(apiPaths, apis[i].Path)
	}
	global.DYCLOUD_DB.Find(&[]system.SysApi{}, "path in (?)", apiPaths).Count(&count)
	if count > 0 {
		return
	}
	err := global.DYCLOUD_DB.Create(&apis).Error
	if err != nil {
		fmt.Println(err)
	}
}

func RegisterMenus(menus ...system.SysBaseMenu) {
	var count int64
	var menuNames []string
	parentMenu := menus[0]
	otherMenus := menus[1:]
	for i := range menus {
		menuNames = append(menuNames, menus[i].Name)
	}
	global.DYCLOUD_DB.Find(&[]system.SysBaseMenu{}, "name in (?)", menuNames).Count(&count)
	if count > 0 {
		return
	}
	err := global.DYCLOUD_DB.Create(&parentMenu).Error
	if err != nil {
		fmt.Println(err)
	}
	for i := range otherMenus {
		pid := parentMenu.ID
		otherMenus[i].ParentId = pid
	}
	err = global.DYCLOUD_DB.Create(&otherMenus).Error
	if err != nil {
		fmt.Println(err)
	}
}
