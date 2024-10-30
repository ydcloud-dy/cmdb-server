package endpoint

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/endpoint"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sEndPointService struct {
	kubernetes.BaseService
}

func (k *K8sEndPointService) GetEndPointList(req endpoint.GetEndPointListReq, uuid uuid.UUID) (*[]corev1.Endpoints, int, error) {
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
	options := metav1.ListOptions{}

	data, err := client.CoreV1().Endpoints(req.Namespace).List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterEndPoints []corev1.Endpoints
	if req.Keyword != "" {
		for _, item := range data.Items {
			if strings.Contains(item.Name, req.Keyword) {
				filterEndPoints = append(filterEndPoints, item)
			}
		}
	} else {
		filterEndPoints = data.Items
	}

	result, err := paginate.Paginate(filterEndPoints, req.Page, req.PageSize)

	return result.(*[]corev1.Endpoints), len(filterEndPoints), nil
}
func (k *K8sEndPointService) DescribeEndPoint(req endpoint.DescribeEndPointReq, uuid uuid.UUID) (*corev1.Endpoints, error) {
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
	EndPointIns, err := client.CoreV1().Endpoints(req.Namespace).Get(context.TODO(), req.EndPointName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return EndPointIns, nil
}
func (k *K8sEndPointService) UpdateEndPoint(req endpoint.UpdateEndPointReq, uuid uuid.UUID) (*corev1.Endpoints, error) {
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
	var EndPointIns *corev1.Endpoints
	err = json.Unmarshal([]byte(tmp), &EndPointIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.CoreV1().Endpoints(req.Namespace).Update(context.TODO(), EndPointIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sEndPointService) DeleteEndPoint(req endpoint.DeleteEndPointReq, uuid uuid.UUID) error {
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

	err = client.CoreV1().Endpoints(req.Namespace).Delete(context.TODO(), req.EndPointName, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())

		return err
	}
	return nil
}
func (k *K8sEndPointService) CreateEndPoint(req endpoint.CreateEndPointReq, uuid uuid.UUID) (*corev1.Endpoints, error) {
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
	var EndPoint *corev1.Endpoints
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &EndPoint)
	ins, err := client.CoreV1().Endpoints(req.Namespace).Create(context.TODO(), EndPoint, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
