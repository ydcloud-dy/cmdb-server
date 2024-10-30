package request

import (
	"DYCLOUD/model/common/request"
	organization "DYCLOUD/plugin/organization/model"
)

type OrganizationSearch struct {
	organization.Organization
	request.PageInfo
}

type OrgUserSearch struct {
	organization.OrgUser
	UserName string `json:"userName" form:"userName"`
	request.PageInfo
}
