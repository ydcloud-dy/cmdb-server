package service

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/service"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sSvcService struct {
	kubernetes.BaseService
}

func (k *K8sSvcService) GetServiceList(req service.GetServiceListReq, uuid uuid.UUID) (*[]v1.Service, int, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	data, err := client.CoreV1().Services(req.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var serviceList []v1.Service
	if req.Keyword != "" {
		for _, service := range data.Items {
			if strings.Contains(service.Name, req.Keyword) {
				serviceList = append(serviceList, service)
			}
		}
	} else {
		serviceList = data.Items
	}

	result, err := paginate.Paginate(serviceList, req.Page, req.PageSize)

	return result.(*[]v1.Service), len(serviceList), nil
}
func (k *K8sSvcService) DescribeService(req service.DescribeServiceReq, uuid uuid.UUID) (*v1.Service, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	service, err := client.CoreV1().Services(req.Namespace).Get(context.TODO(), req.ServiceName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return service, nil
}
func (k *K8sSvcService) UpdateService(req service.UpdateServiceReq, uuid uuid.UUID) (*v1.Service, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	data, err := json.Marshal(req.Content)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	tmp := string(data)
	var serviceIns *v1.Service
	err = json.Unmarshal([]byte(tmp), &serviceIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.CoreV1().Services(req.Namespace).Update(context.TODO(), serviceIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sSvcService) DeleteService(req service.DeleteServiceReq, uuid uuid.UUID) error {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return err
	}

	err = client.CoreV1().Services(req.Namespace).Delete(context.TODO(), req.ServiceName, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())

		return err
	}
	return nil
}
func (k *K8sSvcService) CreateService(req service.CreateServiceReq, uuid uuid.UUID) (*v1.Service, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	data, err := json.Marshal(req.Content)
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	var serviceIns *v1.Service
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &serviceIns)
	ins, err := client.CoreV1().Services(req.Namespace).Create(context.TODO(), serviceIns, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
