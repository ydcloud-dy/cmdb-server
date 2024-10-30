package tencent

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type Region struct {
}

func NewRegion() *Region {
	return &Region{}
}

func (r *Region) List(AccessKeyID, AccessKeySecret string) (list []*cvm.RegionInfo, err error) {
	credential := common.NewCredential(AccessKeyID, AccessKeySecret)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := cvm.NewClient(credential, "", cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := cvm.NewDescribeRegionsRequest()

	// 返回的resp是一个DescribeRegionsResponse的实例，与请求对象对应
	response, err := client.DescribeRegions(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return
	}
	if err != nil {
		panic(err)
	}

	return response.Response.RegionSet, err

}
