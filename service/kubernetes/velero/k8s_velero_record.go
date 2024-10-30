package velero

import (
	"DYCLOUD/global"
	veleroReq "DYCLOUD/model/velero/request"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid/v5"
	v1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	veleroclientset "github.com/vmware-tanzu/velero/pkg/generated/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sVeleroRecordsService struct {
	kubernetes.BaseService
}

func (k8sVeleroRecordsService *K8sVeleroRecordsService) GetK8sVeleroRecordList(
	req veleroReq.K8sVeleroRecordsSearchReq, uuid uuid.UUID) (list *[]v1.Backup, total int, err error) {

	fmt.Print(req)
	kubernetes, err := k8sVeleroRecordsService.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	config, err := kubernetes.Config()
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	veleroClient, err := veleroclientset.NewForConfig(config)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	data, err := veleroClient.VeleroV1().Backups(req.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var backupList []v1.Backup
	if req.Keyword != "" {
		for _, item := range data.Items {
			if strings.Contains(item.Name, req.Keyword) {
				backupList = append(backupList, item)
			}
		}
	} else {
		backupList = data.Items
	}

	result, err := paginate.Paginate(backupList, req.Page, req.PageSize)
	datas, _ := result.(*[]v1.Backup)

	return datas, len(backupList), nil
}
func (K8sVeleroRecordsService *K8sVeleroRecordsService) DescribeVeleroRecord(req *veleroReq.DescribeVeleroRecordReq, uuid uuid.UUID) (schedule *v1.Backup, err error) {
	kubernetes, err := K8sVeleroRecordsService.Generic(req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	config, err := kubernetes.Config()
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	veleroClient, err := veleroclientset.NewForConfig(config)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	veleroRecord, err := veleroClient.VeleroV1().Backups(req.Namespace).Get(context.TODO(), req.VeleroRecordName, metav1.GetOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	return veleroRecord, nil
}
func (K8sVeleroRecordsService *K8sVeleroRecordsService) DeleteK8sVeleroRecord(req *veleroReq.DeleteVeleroRecordReq, uuid uuid.UUID) (err error) {
	fmt.Print(req)
	kubernetes, err := K8sVeleroRecordsService.Generic(req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return err
	}
	config, err := kubernetes.Config()
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return err
	}
	veleroClient, err := veleroclientset.NewForConfig(config)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return err
	}
	err = veleroClient.VeleroV1().Backups(req.Namespace).Delete(context.TODO(), req.VeleroRecordName, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return err
	}
	return nil
}
func (K8sVeleroRecordsService *K8sVeleroRecordsService) CreateK8sVeleroRecord(req *veleroReq.CreateVeleroRecordReq, uuid uuid.UUID) (*v1.Backup, error) {
	kubernetes, err := K8sVeleroRecordsService.Generic(req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	config, err := kubernetes.Config()
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	veleroClient, err := veleroclientset.NewForConfig(config)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	data, err := json.Marshal(req.Content)
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	var backup *v1.Backup
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &backup)
	ins, err := veleroClient.VeleroV1().Backups("velero").Create(context.TODO(), backup, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
func (K8sVeleroRecordsService *K8sVeleroRecordsService) UpdateK8sVeleroRecord(req *veleroReq.UpdateVeleroRecordReq, uuid uuid.UUID) (*v1.Backup, error) {

	kubernetes, err := K8sVeleroRecordsService.Generic(req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	config, err := kubernetes.Config()
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	veleroClient, err := veleroclientset.NewForConfig(config)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	data, err := json.Marshal(req.Content)
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	var backup *v1.Backup
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &backup)
	ins, err := veleroClient.VeleroV1().Backups("velero").Update(context.TODO(), backup, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
