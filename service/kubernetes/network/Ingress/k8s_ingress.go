package Ingress

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/ingress"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sIngressService struct {
	kubernetes.BaseService
}

func (k *K8sIngressService) GetIngressList(req ingress.GetIngressListReq, uuid uuid.UUID) (*[]v1.Ingress, int, error) {
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
	data, err := client.NetworkingV1().Ingresses(req.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterIngresses []v1.Ingress
	if req.Keyword != "" {
		for _, item := range data.Items {
			if strings.Contains(item.Name, req.Keyword) {
				filterIngresses = append(filterIngresses, item)
			}
		}
	} else {
		filterIngresses = data.Items
	}

	result, err := paginate.Paginate(filterIngresses, req.Page, req.PageSize)

	return result.(*[]v1.Ingress), len(filterIngresses), nil
}
func (k *K8sIngressService) DescribeIngress(req ingress.DescribeIngressReq, uuid uuid.UUID) (*v1.Ingress, error) {
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
	IngressIns, err := client.NetworkingV1().Ingresses(req.Namespace).Get(context.TODO(), req.IngressName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return IngressIns, nil
}
func (k *K8sIngressService) UpdateIngress(req ingress.UpdateIngressReq, uuid uuid.UUID) (*v1.Ingress, error) {
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
	var IngressIns *v1.Ingress
	err = json.Unmarshal([]byte(tmp), &IngressIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.NetworkingV1().Ingresses(req.Namespace).Update(context.TODO(), IngressIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sIngressService) DeleteIngress(req ingress.DeleteIngressReq, uuid uuid.UUID) error {
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

	err = client.NetworkingV1().Ingresses(req.Namespace).Delete(context.TODO(), req.IngressName, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())

		return err
	}
	return nil
}
func (k *K8sIngressService) CreateIngress(req ingress.CreateIngressReq, uuid uuid.UUID) (*v1.Ingress, error) {
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
	var Ingress *v1.Ingress
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &Ingress)
	ins, err := client.NetworkingV1().Ingresses(req.Namespace).Create(context.TODO(), Ingress, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
