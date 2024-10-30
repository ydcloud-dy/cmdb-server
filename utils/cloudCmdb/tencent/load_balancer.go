package tencent

import (
	"DYCLOUD/global"
	model "DYCLOUD/model/cloudCmdb"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"go.uber.org/zap"
	"strings"
)

type LoadBalancer struct {
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{}
}

func (l *LoadBalancer) getLoadBalancerIP(ip []*string) string {
	if len(ip) == 0 {
		return ""
	} else {
		return *ip[0]
	}
}

func (l *LoadBalancer) status(status *uint64) string {
	if _, ok := LoadBalancerStatus[*status]; ok {
		return LoadBalancerStatus[*status]
	}

	return ""
}

func (l *LoadBalancer) get(client *clb.Client, pageNumber int64, pageSize int64) ([]*clb.LoadBalancer, error) {
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := clb.NewDescribeLoadBalancersRequest()

	request.Offset = common.Int64Ptr(pageNumber)
	request.Limit = common.Int64Ptr(pageSize)

	// 返回的resp是一个DescribeLoadBalancersResponse的实例，与请求对象对应
	response, err := client.DescribeLoadBalancers(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, err
	}

	return response.Response.LoadBalancerSet, err
}

func (l *LoadBalancer) List(cloudId uint, region model.CloudRegions, AccessKeyID, AccessKeySecret string) (list []model.LoadBalancer, err error) {
	credential := common.NewCredential(AccessKeyID, AccessKeySecret)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "clb.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := clb.NewClient(credential, strings.ReplaceAll(region.RegionId, "tencent-", ""), cpf)
	var pageNumber int64 = 0
	var pageSize int64 = 30

	for {
		response, err := l.get(client, pageNumber, pageSize)
		if err != nil {
			global.DYCLOUD_LOG.Error("clb getInstances fail!", zap.Error(err))
			return list, err
		}

		for _, instance := range response {
			privateAddr := ""
			publicAddr := ""

			if *instance.LoadBalancerType == "INTERNAL" {
				privateAddr = l.getLoadBalancerIP(instance.LoadBalancerVips)
			} else if *instance.LoadBalancerType == "OPEN" {
				publicAddr = l.getLoadBalancerIP(instance.LoadBalancerVips)
			}

			list = append(list, model.LoadBalancer{
				Name:            *instance.LoadBalancerName,
				InstanceId:      *instance.LoadBalancerId,
				PrivateAddr:     privateAddr,
				PublicAddr:      publicAddr,
				Region:          strings.ReplaceAll(region.RegionId, "tencent-", ""),
				RegionName:      region.RegionName,
				Status:          l.status(instance.Status),
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
