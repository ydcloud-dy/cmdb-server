package pvc

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/pvc"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	"golang.org/x/net/context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sPvcService struct {
	kubernetes.BaseService
}

func (k *K8sPvcService) GetPvcList(req pvc.GetPvcListReq, uuid uuid.UUID) (*[]v1.PersistentVolumeClaim, int, error) {
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
	data, err := client.CoreV1().PersistentVolumeClaims(req.Namespace).List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterPVC []v1.PersistentVolumeClaim
	if req.Keyword != "" {
		for _, item := range data.Items {
			if strings.Contains(item.Name, req.Keyword) {
				filterPVC = append(filterPVC, item)
			}
		}
	} else {
		filterPVC = data.Items
	}

	result, err := paginate.Paginate(filterPVC, req.Page, req.PageSize)

	return result.(*[]v1.PersistentVolumeClaim), len(filterPVC), nil

}
func (k *K8sPvcService) DescribePersistentVolumeClaim(req pvc.DescribePVCReq, uuid uuid.UUID) (*v1.PersistentVolumeClaim, error) {
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
	PersistentVolumeClaimIns, err := client.CoreV1().PersistentVolumeClaims(req.Namespace).Get(context.TODO(), req.PVCName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return PersistentVolumeClaimIns, nil
}
func (k *K8sPvcService) UpdatePersistentVolumeClaim(req pvc.UpdatePVCReq, uuid uuid.UUID) (*v1.PersistentVolumeClaim, error) {
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
	var PersistentVolumeClaimIns *v1.PersistentVolumeClaim
	err = json.Unmarshal([]byte(tmp), &PersistentVolumeClaimIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.CoreV1().PersistentVolumeClaims(req.Namespace).Update(context.TODO(), PersistentVolumeClaimIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sPvcService) DeletePersistentVolumeClaim(req pvc.DeletePVCReq, uuid uuid.UUID) error {
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

	err = client.CoreV1().PersistentVolumeClaims(req.Namespace).Delete(context.TODO(), req.PVCName, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())

		return err
	}
	return nil
}
func (k *K8sPvcService) CreatePersistentVolumeClaim(req pvc.CreatePVCReq, uuid uuid.UUID) (*v1.PersistentVolumeClaim, error) {
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
	var namespace *v1.PersistentVolumeClaim
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &namespace)
	ins, err := client.CoreV1().PersistentVolumeClaims(req.Namespace).Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
