package huawei

import (
	"DYCLOUD/global"
	model "DYCLOUD/model/cloudCmdb"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	huaweimodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	reg "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/region"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type ECS struct {
}

func NewECS() *ECS {
	return &ECS{}
}

func (e *ECS) getInstanceIP(address map[string][]huaweimodel.ServerAddress, tp string) string {
	if len(address) == 0 {
		return ""
	}

	addressString := ""
	for k, _ := range address {
		for _, addres := range address[k] {
			if tp == addres.OSEXTIPStype.Value() {
				addressString += addres.Addr + " "
			}
		}
	}

	return addressString
}

func (e *ECS) status(status string) string {
	if _, ok := ECSStatus[status]; ok {
		return ECSStatus[status]
	}

	return ""
}

func (e *ECS) get(client *ecs.EcsClient, pageNumber int, pageSize int) (*[]huaweimodel.ServerDetail, error) {
	request := &huaweimodel.ListServersDetailsRequest{}
	limitRequest := int32(pageSize)
	request.Limit = &limitRequest
	offsetRequest := int32(pageNumber)
	request.Offset = &offsetRequest
	response, err := client.ListServersDetails(request)
	if err != nil {
		return nil, err
	}

	return response.Servers, err
}

func (e *ECS) List(cloudId uint, region model.CloudRegions, AccessKeyID, AccessKeySecret string) (list []model.VirtualMachine, err error) {

	auth := basic.NewCredentialsBuilder().WithAk(AccessKeyID).WithSk(AccessKeySecret).Build()
	client := ecs.NewEcsClient(
		ecs.EcsClientBuilder().
			WithRegion(reg.ValueOf(strings.ReplaceAll(region.RegionId, "huawei-", ""))).
			WithCredential(auth).
			Build())

	pageNumber := 1
	pageSize := 30

	for {

		response, err := e.get(client, pageNumber, pageSize)
		if err != nil {
			global.DYCLOUD_LOG.Error("ecs getInstances fail!", zap.Error(err))
			return list, err
		}

		for _, instance := range *response {
			cpu, _ := strconv.Atoi(instance.Flavor.Vcpus)
			memory, _ := strconv.Atoi(instance.Flavor.Ram)
			list = append(list, model.VirtualMachine{
				Name:            instance.Name,
				InstanceId:      instance.Id,
				UserName:        "root",
				Password:        "changeme",
				Port:            "22",
				CPU:             cpu,
				Memory:          memory,
				OS:              instance.Metadata["image_name"],
				OSType:          instance.Metadata["os_type"],
				PrivateAddr:     e.getInstanceIP(instance.Addresses, "fixed"),
				PublicAddr:      e.getInstanceIP(instance.Addresses, "floating"),
				Region:          strings.ReplaceAll(region.RegionId, "huawei-", ""),
				RegionName:      region.RegionName,
				Status:          e.status(instance.Status),
				CreationTime:    instance.Created,
				CloudPlatformId: cloudId,
			})
		}

		if len(*response) < pageSize {
			break
		}

		pageNumber++
	}

	return list, err
}
