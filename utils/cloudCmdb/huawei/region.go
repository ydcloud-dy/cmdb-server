package huawei

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	iam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"
)

type Region struct {
}

func NewRegion() *Region {
	return &Region{}
}

type Regions struct {
	Id       string
	Name     string
	RegionId string
}

func (r *Region) ListProjects(AccessKeyID, AccessKeySecret string) (rojects *[]model.ProjectResult, err error) {
	auth := global.NewCredentialsBuilder().WithAk(AccessKeyID).WithSk(AccessKeySecret).Build()
	client := iam.NewIamClient(iam.IamClientBuilder().WithRegion(region.ValueOf("cn-east-3")).WithCredential(auth).Build())

	request := &model.KeystoneListProjectsRequest{}
	response, err := client.KeystoneListProjects(request)
	if err != nil {
		return nil, err
	}

	return response.Projects, err
}
func (r *Region) List(AccessKeyID, AccessKeySecret string) (list []Regions, err error) {
	auth := global.NewCredentialsBuilder().WithAk(AccessKeyID).WithSk(AccessKeySecret).Build()
	client := iam.NewIamClient(iam.IamClientBuilder().WithRegion(region.ValueOf("cn-east-3")).WithCredential(auth).Build())
	request := &model.KeystoneListRegionsRequest{}

	response, err := client.KeystoneListRegions(request)
	if err != nil {
		return nil, err
	}

	var regions []Regions
	for _, re := range *response.Regions {
		projects, err := r.ListProjects(AccessKeyID, AccessKeySecret)
		if err != nil {
			return nil, err
		}

		for _, project := range *projects {
			if re.Id == project.Name {
				regions = append(regions, Regions{
					Id:       re.Id,
					Name:     re.Locales.ZhCn,
					RegionId: re.Id,
				})
			}
		}
	}

	return regions, err
}
