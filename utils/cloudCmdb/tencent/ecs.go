package tencent

import (
	"DYCLOUD/global"
	model "DYCLOUD/model/cloudCmdb"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"go.uber.org/zap"
	"strings"
)

type ECS struct {
}

func NewECS() *ECS {
	return &ECS{}
}

func (e *ECS) getInstanceIP(ip []*string) string {
	if len(ip) == 0 {
		return ""
	} else {
		return *ip[0]
	}
}

func (e *ECS) get(client *cvm.Client, pageNumber int64, pageSize int64) (*cvm.DescribeInstancesResponseParams, error) {
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := cvm.NewDescribeInstancesRequest()

	// 返回的resp是一个DescribeInstancesResponse的实例，与请求对象对应
	response, err := client.DescribeInstances(request)

	request.Offset = common.Int64Ptr(pageNumber)
	request.Limit = common.Int64Ptr(pageSize)

	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return response.Response, nil
}

func (e *ECS) status(status string) string {
	if _, ok := ECSStatus[status]; ok {
		return ECSStatus[status]
	}

	return ""
}

func (e *ECS) List(cloudId uint, region model.CloudRegions, AccessKeyID, AccessKeySecret string) (list []model.VirtualMachine, err error) {
	credential := common.NewCredential(AccessKeyID, AccessKeySecret)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := cvm.NewClient(credential, strings.ReplaceAll(region.RegionId, "tencent-", ""), cpf)
	var pageNumber int64 = 0
	var pageSize int64 = 30

	for {
		response, err := e.get(client, pageNumber, pageSize)
		if err != nil {
			global.DYCLOUD_LOG.Error("ecs getInstances fail!", zap.Error(err))
			return list, err
		}

		for _, instance := range response.InstanceSet {
			list = append(list, model.VirtualMachine{
				Name:            *instance.InstanceName,
				InstanceId:      *instance.InstanceId,
				UserName:        "root",
				Password:        "changeme",
				Port:            "22",
				CPU:             int(*instance.CPU),
				Memory:          int(*instance.Memory),
				OS:              *instance.OsName,
				OSType:          *instance.OsName,
				PrivateAddr:     e.getInstanceIP(instance.PrivateIpAddresses),
				PublicAddr:      e.getInstanceIP(instance.PublicIpAddresses),
				Region:          strings.ReplaceAll(region.RegionId, "tencent-", ""),
				RegionName:      region.RegionName,
				Status:          e.status(*instance.InstanceState),
				CreationTime:    *instance.CreatedTime,
				ExpiredTime:     *instance.ExpiredTime,
				CloudPlatformId: cloudId,
			})
		}

		if len(response.InstanceSet) < int(pageSize) {
			break
		}

		pageNumber++
	}

	return list, err

}
