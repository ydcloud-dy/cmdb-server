package utils

import "DYCLOUD/utils"

var (
	CloudVerify = utils.Rules{"Name": {utils.NotEmpty()}, "AccessKeyId": {utils.NotEmpty()}, "AccessKeySecret": {utils.NotEmpty()}, "Platform": {utils.NotEmpty()}}
)
