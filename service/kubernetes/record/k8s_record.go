package record

import (
	"DYCLOUD/model/record"
	"DYCLOUD/utils/kubernetes"
)

type K8sRecordService struct {
	kubernetes.BaseService
}

func (k *K8sRecordService) GetRecordList(req record.GetRecordListReq) (*string, int, error) {
	return nil, 0, nil
}
func (k *K8sRecordService) DescribeRecord(req record.DescribeRecordReq) (*string, error) {
	return nil, nil
}
func (k *K8sRecordService) UpdateRecord(req record.UpdateRecordReq) (*string, error) {
	return nil, nil

}
func (k *K8sRecordService) CreateRecord(req record.CreateRecordReq) (*string, error) {
	return nil, nil

}
func (k *K8sRecordService) DeleteRecord(req record.DeleteRecordReq) error {
	return nil
}
