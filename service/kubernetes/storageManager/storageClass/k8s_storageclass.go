package storageClass

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/storageClass"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sStorageClassService struct {
	kubernetes.BaseService
}

func (k *K8sStorageClassService) GetStorageClassList(req storageClass.GetStorageClassListReq, uuid uuid.UUID) (*[]v1.StorageClass, int, error) {
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

	data, err := client.StorageV1().StorageClasses().List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterPod []v1.StorageClass
	if req.Keyword != "" {
		for _, storageClass := range data.Items {
			if strings.Contains(storageClass.Name, req.Keyword) {
				filterPod = append(filterPod, storageClass)
			}
		}
	} else {
		filterPod = data.Items
	}

	result, err := paginate.Paginate(filterPod, req.Page, req.PageSize)

	return result.(*[]v1.StorageClass), len(filterPod), nil
}
func (k *K8sStorageClassService) DescribeStorageClass(req storageClass.DescribeStorageClassReq, uuid uuid.UUID) (*v1.StorageClass, error) {
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
	StorageClass, err := client.StorageV1().StorageClasses().Get(context.TODO(), req.StorageClassName, metav1.GetOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	return StorageClass, nil
}
func (k *K8sStorageClassService) UpdateStorageClass(req storageClass.UpdateStorageClassReq, uuid uuid.UUID) (*v1.StorageClass, error) {
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
	var StorageClassIns *v1.StorageClass
	err = json.Unmarshal([]byte(tmp), &StorageClassIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.StorageV1().StorageClasses().Update(context.TODO(), StorageClassIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sStorageClassService) CreateStorageClass(req storageClass.CreateStorageClassReq, uuid uuid.UUID) (*v1.StorageClass, error) {
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
	var StorageClass *v1.StorageClass
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &StorageClass)
	ins, err := client.StorageV1().StorageClasses().Create(context.TODO(), StorageClass, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
func (k *K8sStorageClassService) DeleteStorageClass(req storageClass.DeleteStorageClassReq, uuid uuid.UUID) error {
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
	err = client.StorageV1().StorageClasses().Delete(context.TODO(), req.StorageClassName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
