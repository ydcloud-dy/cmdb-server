package huawei

import (
	"DYCLOUD/global"
	model "DYCLOUD/model/cloudCmdb"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	reg "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v2/region"
	elb "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3"
	huaweimodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/model"
	"go.uber.org/zap"
	"strings"
)

type LoadBalancer struct {
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{}
}

func (l *LoadBalancer) getLoadBalancerEips(address []huaweimodel.PublicIpInfo) string {
	addString := ""

	if len(address) <= 0 {
		return ""
	}

	for _, addres := range address {
		addString += addres.PublicipAddress + " "
	}

	return addString
}

func (l *LoadBalancer) status(status string) string {
	if _, ok := LoadBalancerStatus[status]; ok {
		return LoadBalancerStatus[status]
	}

	return ""
}

func (l *LoadBalancer) get(client *elb.ElbClient, marker string, pageSize int) (*[]huaweimodel.LoadBalancer, error) {
	request := &huaweimodel.ListLoadBalancersRequest{}
	limitRequest := int32(pageSize)
	request.Limit = &limitRequest

	if marker != "" {
		markerRequest := marker
		request.Marker = &markerRequest
		pageReverseRequest := false
		request.PageReverse = &pageReverseRequest
	}

	response, err := client.ListLoadBalancers(request)

	if err != nil {
		return nil, err
	}

	return response.Loadbalancers, nil
}

func (l *LoadBalancer) List(cloudId uint, region model.CloudRegions, AccessKeyID, AccessKeySecret string) (list []model.LoadBalancer, err error) {
	auth := basic.NewCredentialsBuilder().WithAk(AccessKeyID).WithSk(AccessKeySecret).Build()
	client := elb.NewElbClient(elb.ElbClientBuilder().
		WithRegion(reg.ValueOf(strings.ReplaceAll(region.RegionId, "huawei-", ""))).
		WithCredential(auth).
		Build())

	pageNumber := 1
	pageSize := 30
	marker := ""

	for {

		response, err := l.get(client, marker, pageSize)
		if err != nil {
			global.DYCLOUD_LOG.Error("ecs getInstances fail!", zap.Error(err))
			return list, err
		}

		for _, instance := range *response {
			list = append(list, model.LoadBalancer{
				Name:            instance.Name,
				InstanceId:      instance.Id,
				PrivateAddr:     instance.VipAddress,
				PublicAddr:      l.getLoadBalancerEips(instance.Publicips),
				Region:          strings.ReplaceAll(region.RegionId, "huawei-", ""),
				RegionName:      region.RegionName,
				Status:          l.status(instance.ProvisioningStatus),
				CreationTime:    instance.CreatedAt,
				CloudPlatformId: cloudId,
			})

			marker = instance.Id
			if len(*response) < pageSize {
				marker = ""
			}
		}

		if len(*response) < pageSize {
			break
		}

		pageNumber++
	}
	return list, err
}
