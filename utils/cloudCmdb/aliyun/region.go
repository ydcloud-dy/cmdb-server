package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

type Region struct {
}

func NewRegion() *Region {
	return &Region{}
}

func (r *Region) List(AccessKeyID, AccessKeySecret string) (list []ecs.Region, err error) {
	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(AccessKeyID, AccessKeySecret)
	client, err := ecs.NewClientWithOptions("cn-hangzhou", config, credential)
	if err != nil {
		return list, err
	}

	request := ecs.CreateDescribeRegionsRequest()
	request.Scheme = "https"
	response, err := client.DescribeRegions(request)
	if err != nil {
		return list, err
	}

	return response.Regions.Region, err
}
