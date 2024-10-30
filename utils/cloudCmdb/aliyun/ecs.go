package aliyun

import (
	"DYCLOUD/global"
	model "DYCLOUD/model/cloudCmdb"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"go.uber.org/zap"
	"strings"
)

type ECS struct {
}

func NewECS() *ECS {
	return &ECS{}
}

func (e *ECS) get(client *ecs.Client, pageNumber int, pageSize int) (*ecs.DescribeInstancesResponse, error) {
	request := ecs.CreateDescribeInstancesRequest()
	request.PageNumber = requests.Integer(fmt.Sprintf("%d", pageNumber))
	request.PageSize = requests.Integer(fmt.Sprintf("%d", pageSize))

	response, err := client.DescribeInstances(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (e *ECS) getInstanceIP(ip []string) string {
	if len(ip) == 0 {
		return ""
	} else {
		return ip[0]
	}
}

func (e *ECS) status(status string) string {
	if _, ok := ECSStatus[status]; ok {
		return ECSStatus[status]
	}

	return ""
}

func (e *ECS) List(cloudId uint, region model.CloudRegions, AccessKeyID, AccessKeySecret string) (list []model.VirtualMachine, err error) {

	client, err := ecs.NewClientWithAccessKey(strings.ReplaceAll(region.RegionId, "aliyun-", ""), AccessKeyID, AccessKeySecret)
	if err != nil {
		global.DYCLOUD_LOG.Error("ecs new Client fail!", zap.Error(err))
		return
	}

	pageNumber := 1
	pageSize := 30

	for {
		response, err := e.get(client, pageNumber, pageSize)
		if err != nil {
			global.DYCLOUD_LOG.Error("ecs getInstances fail!", zap.Error(err))
			return list, err
		}

		for _, instance := range response.Instances.Instance {
			list = append(list, model.VirtualMachine{
				Name:            instance.InstanceName,
				InstanceId:      instance.InstanceId,
				UserName:        "root",
				Password:        "changeme",
				Port:            "22",
				CPU:             instance.Cpu,
				Memory:          instance.Memory,
				OS:              instance.OSName,
				OSType:          instance.OSType,
				PrivateAddr:     e.getInstanceIP(instance.VpcAttributes.PrivateIpAddress.IpAddress),
				PublicAddr:      e.getInstanceIP(instance.PublicIpAddress.IpAddress),
				Region:          strings.ReplaceAll(region.RegionId, "aliyun-", ""),
				RegionName:      region.RegionName,
				Status:          e.status(instance.Status),
				CreationTime:    instance.CreationTime,
				ExpiredTime:     instance.ExpiredTime,
				CloudPlatformId: cloudId,
			})
		}

		if len(response.Instances.Instance) < pageSize {
			break
		}

		pageNumber++
	}

	return list, err
}
