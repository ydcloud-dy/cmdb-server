package tencent

import (
	"DYCLOUD/global"
	model "DYCLOUD/model/cloudCmdb"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"go.uber.org/zap"
	"strings"
)

type RDS struct {
}

func NewRDS() *RDS {
	return &RDS{}
}

func (r *RDS) status(status int64) string {
	if _, ok := RdsStatus[status]; ok {
		return RdsStatus[status]
	}

	return ""
}

func (r *RDS) get(client *cdb.Client, pageNumber int64, pageSize int64) ([]*cdb.InstanceInfo, error) {
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := cdb.NewDescribeDBInstancesRequest()

	request.Offset = common.Uint64Ptr(0)
	request.Limit = common.Uint64Ptr(30)

	// 返回的resp是一个DescribeDBInstancesResponse的实例，与请求对象对应
	response, err := client.DescribeDBInstances(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return response.Response.Items, err
}

func (r *RDS) List(cloudId uint, region model.CloudRegions, AccessKeyID, AccessKeySecret string) (list []model.RDS, err error) {
	credential := common.NewCredential(AccessKeyID, AccessKeySecret)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cdb.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := cdb.NewClient(credential, strings.ReplaceAll(region.RegionId, "tencent-", ""), cpf)

	var pageNumber int64 = 0
	var pageSize int64 = 30

	for {
		response, err := r.get(client, pageNumber, pageSize)
		if err != nil {
			global.DYCLOUD_LOG.Error("clb getInstances fail!", zap.Error(err))
			return list, err
		}

		for _, instance := range response {
			list = append(list, model.RDS{
				Name:            *instance.InstanceName,
				InstanceId:      *instance.InstanceId,
				Region:          strings.ReplaceAll(region.RegionId, "tencent-", ""),
				RegionName:      region.RegionName,
				Status:          r.status(*instance.Status),
				CreationTime:    *instance.CreateTime,
				CloudPlatformId: cloudId,
			})
		}

		if len(response) < int(pageSize) {
			break
		}

		pageNumber++
	}

	return list, err

}
