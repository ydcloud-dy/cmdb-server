package velero

import (
	"DYCLOUD/global"
	veleroReq "DYCLOUD/model/velero/request"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	veleroclientset "github.com/vmware-tanzu/velero/pkg/generated/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sVeleroTasksService struct {
	kubernetes.BaseService
}

// CreateK8sVeleroTasks 创建k8sVeleroTasks表记录
func (k8sVeleroTasksService *K8sVeleroTasksService) CreateK8sVeleroTasks(req *veleroReq.CreateVeleroTaskReq, uuid uuid.UUID) (*v1.Schedule, error) {
	kubernetes, err := k8sVeleroTasksService.Generic(req, uuid)
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
	var Schedule *v1.Schedule
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &Schedule)
	ins, err := veleroClient.VeleroV1().Schedules("velero").Create(context.TODO(), Schedule, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}

// DeleteK8sVeleroTasks 删除k8sVeleroTasks表记录
func (k8sVeleroTasksService *K8sVeleroTasksService) DeleteK8sVeleroTasks(req *veleroReq.DeleteVeleroTaskReq, uuid uuid.UUID) (err error) {
	kubernetes, err := k8sVeleroTasksService.Generic(req, uuid)
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
	err = veleroClient.VeleroV1().Schedules(req.Namespace).Delete(context.TODO(), req.VeleroTaskName, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return err
	}
	return nil
}

// UpdateK8sVeleroTasks 更新k8sVeleroTasks表记录
// Author [piexlmax](https://github.com/piexlmax)
func (k8sVeleroTasksService *K8sVeleroTasksService) UpdateK8sVeleroTasks(req *veleroReq.UpdateVeleroTaskReq, uuid uuid.UUID) (*v1.Schedule, error) {

	kubernetes, err := k8sVeleroTasksService.Generic(req, uuid)
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
	var Schedule *v1.Schedule
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &Schedule)
	ins, err := veleroClient.VeleroV1().Schedules("velero").Update(context.TODO(), Schedule, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}

// GetK8sVeleroTasks 根据ID获取k8sVeleroTasks表记录
func (k8sVeleroTasksService *K8sVeleroTasksService) DescribeVeleroTask(req *veleroReq.DescribeVeleroTaskReq, uuid uuid.UUID) (schedule *v1.Schedule, err error) {
	kubernetes, err := k8sVeleroTasksService.Generic(req, uuid)
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
	veleroTask, err := veleroClient.VeleroV1().Schedules(req.Namespace).Get(context.TODO(), req.VeleroTaskName, metav1.GetOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	return veleroTask, nil
}

// GetK8sVeleroTasksInfoList 分页获取k8sVeleroTasks表记录
// Author [piexlmax](https://github.com/piexlmax)
func (k8sVeleroTasksService *K8sVeleroTasksService) GetK8sVeleroTasksInfoList(
	req veleroReq.K8sVeleroTasksSearchReq, uuid uuid.UUID) (list *[]v1.Schedule, total int, err error) {

	kubernetes, err := k8sVeleroTasksService.Generic(&req, uuid)
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
	data, err := veleroClient.VeleroV1().Schedules(req.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var scheduleList []v1.Schedule
	if req.Keyword != "" {
		for _, item := range data.Items {
			if strings.Contains(item.Name, req.Keyword) {
				scheduleList = append(scheduleList, item)
			}
		}
	} else {
		scheduleList = data.Items
	}

	result, err := paginate.Paginate(scheduleList, req.Page, req.PageSize)
	datas, _ := result.(*[]v1.Schedule)

	return datas, len(scheduleList), nil
}
