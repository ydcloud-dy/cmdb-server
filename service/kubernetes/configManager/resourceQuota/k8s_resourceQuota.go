package resourceQuota

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/resourceQuota"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sResourceQuotaService struct {
	kubernetes.BaseService
}

func (k *K8sResourceQuotaService) GetResourceQuotaList(req resourceQuota.GetResourceQuotaListReq, uuid uuid.UUID) (*[]v1.ResourceQuota, int, error) {
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

	options := metav1.ListOptions{LabelSelector: req.LabelSelector}
	data, err := client.CoreV1().ResourceQuotas(req.Namespace).List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterResourceQuota []v1.ResourceQuota
	if req.Keyword != "" {
		for _, item := range data.Items {
			if strings.Contains(item.Name, req.Keyword) {
				filterResourceQuota = append(filterResourceQuota, item)
			}
		}
	} else {
		filterResourceQuota = data.Items
	}

	result, err := paginate.Paginate(filterResourceQuota, req.Page, req.PageSize)

	return result.(*[]v1.ResourceQuota), len(filterResourceQuota), nil
}
func (k *K8sResourceQuotaService) DescribeResourceQuota(req resourceQuota.DescribeResourceQuotaReq, uuid uuid.UUID) (*v1.ResourceQuota, error) {
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
	resourceQuotaIns, err := client.CoreV1().ResourceQuotas(req.Namespace).Get(context.TODO(), req.ResourceQuotaName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return resourceQuotaIns, nil
}
func (k *K8sResourceQuotaService) UpdateResourceQuota(req resourceQuota.UpdateResourceQuotaReq, uuid uuid.UUID) (*v1.ResourceQuota, error) {
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
	var resourceQuotaIns *v1.ResourceQuota
	err = json.Unmarshal([]byte(tmp), &resourceQuotaIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.CoreV1().ResourceQuotas(req.Namespace).Update(context.TODO(), resourceQuotaIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sResourceQuotaService) DeleteResourceQuota(req resourceQuota.DeleteResourceQuotaReq, uuid uuid.UUID) error {
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

	err = client.CoreV1().ResourceQuotas(req.Namespace).Delete(context.TODO(), req.ResourceQuotaName, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())

		return err
	}
	return nil
}
func (k *K8sResourceQuotaService) CreateResourceQuota(req resourceQuota.CreateResourceQuotaReq, uuid uuid.UUID) (*v1.ResourceQuota, error) {
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
	var namespace *v1.ResourceQuota
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &namespace)
	ins, err := client.CoreV1().ResourceQuotas(req.Namespace).Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
