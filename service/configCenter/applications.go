package configCenter

import (
	"DYCLOUD/model/configCenter"
	"DYCLOUD/model/configCenter/request"
)

type ApplicationsService struct{}

// GetApplicationsList
//
//	@Description: 获取应用环境列表
//	@receiver e
//	@param req
//	@return envList
//	@return err
func (e *ApplicationsService) GetApplicationsList(req *request.EnvRequest) (envList *[]configCenter.Applications, total int64, err error) {
	return nil, 0, err
}
func (e *ApplicationsService) DescribeApplications(id int) (envList *configCenter.Applications, err error) {
	return nil, err
}
func (e *ApplicationsService) CreateApplications(req *configCenter.Applications) error {
	return nil
}
func (e *ApplicationsService) UpdateApplications(req *configCenter.Applications) (data *configCenter.Applications, err error) {
	return nil, err
}
func (e *ApplicationsService) DeleteApplications(id int) error {
	return nil
}

func (e *ApplicationsService) DeleteApplicationsByIds(ids *request.DeleteEnvByIds) error {
	return nil
}
