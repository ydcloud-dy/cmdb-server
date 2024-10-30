package huawei

import (
	"DYCLOUD/global"
	model "DYCLOUD/model/cloudCmdb"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	rds "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3"
	huaweimodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/model"
	reg "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/region"
	"go.uber.org/zap"
	"strings"
)

type RDS struct {
}

func NewRDS() *RDS {
	return &RDS{}
}

func (r *RDS) getInstanceIP(address []string) string {
	addString := ""

	if len(address) <= 0 {
		return ""
	}

	for _, addres := range address {
		addString += addres + " "
	}

	return addString
}

func (r *RDS) status(status string) string {
	if _, ok := RdsStatus[status]; ok {
		return RdsStatus[status]
	}

	return ""
}

func (r *RDS) get(client *rds.RdsClient, pageNumber int, pageSize int) (*[]huaweimodel.InstanceResponse, error) {
	request := &huaweimodel.ListInstancesRequest{}
	offsetRequest := int32(0)
	request.Offset = &offsetRequest
	limitRequest := int32(30)
	request.Limit = &limitRequest
	response, err := client.ListInstances(request)
	if err != nil {
		return nil, err
	}

	return response.Instances, err
}

func (r *RDS) List(cloudId uint, region model.CloudRegions, AccessKeyID, AccessKeySecret string) (list []model.RDS, err error) {
	auth := basic.NewCredentialsBuilder().WithAk(AccessKeyID).WithSk(AccessKeySecret).Build()
	client := rds.NewRdsClient(
		rds.RdsClientBuilder().
			WithRegion(reg.ValueOf(strings.ReplaceAll(region.RegionId, "huawei-", ""))).
			WithCredential(auth).
			Build())

	pageNumber := 0
	pageSize := 30

	for {

		response, err := r.get(client, pageNumber, pageSize)
		if err != nil {
			global.DYCLOUD_LOG.Error("rds getInstances fail!", zap.Error(err))
			break
		}

		for _, instance := range *response {
			list = append(list, model.RDS{
				Name:            instance.Name,
				InstanceId:      instance.Id,
				PrivateAddr:     r.getInstanceIP(instance.PrivateIps),
				PublicAddr:      r.getInstanceIP(instance.PublicIps),
				Region:          strings.ReplaceAll(region.RegionId, "huawei-", ""),
				RegionName:      region.RegionName,
				Status:          r.status(instance.Status),
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
