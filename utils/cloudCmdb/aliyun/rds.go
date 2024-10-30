package aliyun

import (
	"DYCLOUD/global"
	model "DYCLOUD/model/cloudCmdb"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"go.uber.org/zap"
	"strings"
)

type RDS struct {
}

func NewRDS() *RDS {
	return &RDS{}
}

func (r *RDS) status(status string) string {
	if _, ok := RdsStatus[status]; ok {
		return RdsStatus[status]
	}

	return ""
}

func (r *RDS) get(client *rds.Client, pageNumber int, pageSize int) ([]rds.DBInstance, error) {
	request := rds.CreateDescribeDBInstancesRequest()
	request.Scheme = "https"
	request.PageSize = requests.NewInteger(pageSize)
	request.PageNumber = requests.NewInteger(pageNumber)
	response, err := client.DescribeDBInstances(request)
	if err != nil {
		return nil, err
	}

	return response.Items.DBInstance, err
}

func (r *RDS) List(cloudId uint, region model.CloudRegions, AccessKeyID, AccessKeySecret string) (list []model.RDS, err error) {
	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(AccessKeyID, AccessKeySecret)
	client, err := rds.NewClientWithOptions(strings.ReplaceAll(region.RegionId, "aliyun-", ""), config, credential)
	if err != nil {
		global.DYCLOUD_LOG.Error("Aliyun new Client fail!", zap.Error(err))
		return
	}

	pageNumber := 1
	pageSize := 30

	for {
		response, err := r.get(client, pageNumber, pageSize)
		if err != nil {
			global.DYCLOUD_LOG.Error("RDS getInstances fail!", zap.Error(err))
			return list, err
		}

		for _, instance := range response {
			list = append(list, model.RDS{
				Name:            instance.DBInstanceDescription,
				InstanceId:      instance.DBInstanceId,
				PrivateAddr:     instance.ConnectionString,
				Region:          strings.ReplaceAll(region.RegionId, "aliyun-", ""),
				RegionName:      region.RegionName,
				Status:          r.status(instance.DBInstanceStatus),
				CreationTime:    instance.CreateTime,
				CloudPlatformId: cloudId,
			})
		}

		if len(response) < pageSize {
			break
		}

		pageNumber++
	}

	return list, err
}
