package pv

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/pv"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sPVService struct {
	kubernetes.BaseService
}

func (k *K8sPVService) GetPVList(req pv.GetPVListReq, uuid uuid.UUID) (*[]v1.PersistentVolume, int, error) {
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
	data, err := client.CoreV1().PersistentVolumes().List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterPV []v1.PersistentVolume
	if req.Keyword != "" {
		for _, item := range data.Items {
			if strings.Contains(item.Name, req.Keyword) {
				filterPV = append(filterPV, item)
			}
		}
	} else {
		filterPV = data.Items
	}

	result, err := paginate.Paginate(filterPV, req.Page, req.PageSize)

	return result.(*[]v1.PersistentVolume), len(filterPV), nil
}
func (k *K8sPVService) DescribePV(req pv.DescribePVReq, uuid uuid.UUID) (*v1.PersistentVolume, error) {
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
	PVIns, err := client.CoreV1().PersistentVolumes().Get(context.TODO(), req.PVName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return PVIns, nil
}
func (k *K8sPVService) UpdatePV(req pv.UpdatePVReq, uuid uuid.UUID) (*v1.PersistentVolume, error) {
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
	var PVIns *v1.PersistentVolume
	err = json.Unmarshal([]byte(tmp), &PVIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.CoreV1().PersistentVolumes().Update(context.TODO(), PVIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sPVService) DeletePV(req pv.DeletePVReq, uuid uuid.UUID) error {
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

	err = client.CoreV1().PersistentVolumes().Delete(context.TODO(), req.PVName, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())

		return err
	}
	return nil
}
func (k *K8sPVService) CreatePV(req pv.CreatePVReq, uuid uuid.UUID) (*v1.PersistentVolume, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	data, err := json.Marshal(req.Content)
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	var PV *v1.PersistentVolume
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &PV)
	ins, err := client.CoreV1().PersistentVolumes().Create(context.TODO(), PV, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
