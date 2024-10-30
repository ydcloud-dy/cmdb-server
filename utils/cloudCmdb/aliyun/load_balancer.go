package aliyun

import (
	"DYCLOUD/global"
	model "DYCLOUD/model/cloudCmdb"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type LoadBalancer struct {
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{}
}

func (l *LoadBalancer) get(client *slb.Client, pageNumber int, pageSize int) ([]slb.LoadBalancer, error) {

	request := slb.CreateDescribeLoadBalancersRequest()
	request.Scheme = "https"
	request.PageSize = requests.NewInteger(pageSize)
	request.PageNumber = requests.NewInteger(pageNumber)

	response, err := client.DescribeLoadBalancers(request)
	if err != nil {
		return nil, err
	}

	return response.LoadBalancers.LoadBalancer, err
}

func (l *LoadBalancer) status(status string) string {
	if _, ok := LoadBalancerStatus[status]; ok {
		return LoadBalancerStatus[status]
	}

	return ""
}

func (l *LoadBalancer) getInstanceIP(instance slb.LoadBalancer, tp string) string {
	if instance.AddressType == "intranet" && instance.AddressType == tp {
		return instance.Address
	}

	if instance.AddressType == "internet" && instance.AddressType == tp {
		return instance.Address
	}

	return ""
}

func (l LoadBalancer) List(cloudId uint, region model.CloudRegions, AccessKeyID, AccessKeySecret string) (list []model.LoadBalancer, err error) {
	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(AccessKeyID, AccessKeySecret)
	client, err := slb.NewClientWithOptions(strings.ReplaceAll(region.RegionId, "aliyun-", ""), config, credential)
	if err != nil {
		global.DYCLOUD_LOG.Error("LoadBalancer new Client fail!", zap.Error(err))
		return
	}

	pageNumber := 1
	pageSize := 30

	for {
		response, err := l.get(client, pageNumber, pageSize)
		if err != nil {
			global.DYCLOUD_LOG.Error("LoadBalancer getInstances fail!", zap.Error(err))
			return list, err
		}

		for _, instance := range response {

			bandwidth := ""
			if instance.Bandwidth != 0 {
				bandwidth = strconv.Itoa(instance.Bandwidth)
			}

			list = append(list, model.LoadBalancer{
				Name:            instance.LoadBalancerName,
				InstanceId:      instance.LoadBalancerId,
				PrivateAddr:     l.getInstanceIP(instance, "intranet"),
				PublicAddr:      l.getInstanceIP(instance, "internet"),
				Bandwidth:       bandwidth,
				Region:          strings.ReplaceAll(region.RegionId, "aliyun-", ""),
				RegionName:      region.RegionName,
				Status:          l.status(instance.LoadBalancerStatus),
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
