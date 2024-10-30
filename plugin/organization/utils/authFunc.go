package utils

import (
	"DYCLOUD/global"
	"DYCLOUD/model/system"
	"DYCLOUD/plugin/organization/model"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
)

// uint 去重方法
func Uniq(array []uint) []uint {
	var uintMap = make(map[uint]bool)
	var res []uint
	for _, u := range array {
		if !uintMap[u] {
			uintMap[u] = true
			res = append(res, u)
		}
	}
	return res
}

const (
	Node    = 0 // 无资源权限
	Self    = 1 // 仅自己
	Current = 2 // 当前部门
	Deep    = 3 // 当前部门及以下
	All     = 4 // 所有
)

// 获取当前部门ID
func GetSelfOrg(id uint) []uint {
	var orgUser []model.OrgUser
	err := global.DYCLOUD_DB.Find(&orgUser, "sys_user_id = ?", id).Error
	if err != nil {
		return []uint{}
	}
	var orgId []uint
	for _, m := range orgUser {
		orgId = append(orgId, m.OrganizationID)
	}
	return Uniq(orgId)
}

// 获取所有部门
func GetAllOrg() []model.Organization {
	var orgUser []model.Organization
	err := global.DYCLOUD_DB.Find(&orgUser).Error
	if err != nil {
		return []model.Organization{}
	}
	return orgUser
}

// 获取所有部门ID
func GetAllOrgID() []uint {
	orgUser := GetAllOrg()
	if len(orgUser) == 0 {
		return []uint{}
	}
	var orgids []uint
	for _, organization := range orgUser {
		orgids = append(orgids, organization.ID)
	}
	return Uniq(orgids)
}

// 获取当前部门及以下部门id
func GetDeepOrg(id uint) []uint {
	orgId := GetSelfOrg(id)
	if len(orgId) == 0 {
		return []uint{}
	}
	orgs := GetAllOrg()
	if len(orgs) == 0 {
		return []uint{}
	}
	orgids := findChildren(orgId, orgs)
	return Uniq(append(orgids, orgId...))
}

// 获取当前部门及以下部门的递归方法
func findChildren(ids []uint, orgs []model.Organization) []uint {
	var idsMap = make(map[uint]bool)
	var resIDs []uint
	for _, id := range ids {
		idsMap[id] = true
	}
	for _, org := range orgs {
		if idsMap[org.ParentID] {
			resIDs = append(resIDs, org.ID)
		}
	}
	if len(resIDs) == 0 {
		resIDs = append(resIDs, ids...)
		return resIDs
	}
	dids := findChildren(resIDs, orgs)
	resIDs = append(resIDs, ids...)
	resIDs = append(resIDs, dids...)
	return resIDs
}

// 获取当前部门的用户id
func GetCurrentUserIDs(id uint) []uint {
	orgId := GetSelfOrg(id)
	if len(orgId) == 0 {
		return []uint{}
	}
	return GetUsersByOrgIds(orgId)
}

// 获取当前部门及以下的用户id
func GetDeepUserIDs(id uint) []uint {
	orgids := GetDeepOrg(id)
	if len(orgids) == 0 {
		return []uint{}
	}
	return GetUsersByOrgIds(orgids)
}

// 根据部门获取部门下用户ID
func GetUsersByOrgIds(orgIds []uint) []uint {
	var orgUser []model.OrgUser
	err := global.DYCLOUD_DB.Find(&orgUser, "organization_id in (?)", orgIds).Error
	if err != nil {
		return []uint{}
	}
	var userIDS []uint
	for _, m := range orgUser {
		userIDS = append(userIDS, m.SysUserID)
	}
	return Uniq(userIDS)
}

// 获取所有用户ID
func GetAllUserIDs() []uint {
	var users []system.SysUser
	err := global.DYCLOUD_DB.Find(&users).Error
	if err != nil {
		return []uint{}
	}
	var usersID []uint
	for _, sysUser := range users {
		usersID = append(usersID, sysUser.ID)
	}
	return Uniq(usersID)
}

// 自动获取当前用户拥有的权限的用户ID
func GetUserIDS(c *gin.Context) []uint {
	user := utils.GetUserInfo(c)
	var data model.DataAuthority
	err := global.DYCLOUD_DB.First(&data, "authority_id = ?", user.AuthorityId).Error
	if err != nil {
		return []uint{}
	}
	switch data.AuthorityType {
	case Node:
		return []uint{}
	case Self:
		return []uint{user.BaseClaims.ID}
	case Current:
		return GetCurrentUserIDs(user.BaseClaims.ID)
	case Deep:
		return GetDeepUserIDs(user.BaseClaims.ID)
	case All:
		return GetAllUserIDs()
	}
	return []uint{}
}

// 自动获取当前用户拥有的权限的部门ID
func GetOrgIDS(c *gin.Context) []uint {
	user := utils.GetUserInfo(c)
	var data model.DataAuthority
	err := global.DYCLOUD_DB.First(&data, "authority_id = ?", user.AuthorityId).Error
	if err != nil {
		return []uint{}
	}
	switch data.AuthorityType {
	case Node:
		return []uint{}
	case Self:
		return GetSelfOrg(user.BaseClaims.ID)
	case Current:
		return GetSelfOrg(user.BaseClaims.ID)
	case Deep:
		return GetDeepOrg(user.BaseClaims.ID)
	case All:
		return GetAllOrgID()
	}
	return []uint{}
}
