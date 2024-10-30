package configmap

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/configmap"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sConfigMapService struct {
	kubernetes.BaseService
}

func (k *K8sConfigMapService) GetConfigMapList(req configmap.GetConfigMapListReq, uuid uuid.UUID) (*[]v1.ConfigMap, int, error) {
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
	data, err := client.CoreV1().ConfigMaps(req.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var configMapList []v1.ConfigMap
	if req.Keyword != "" {
		for _, item := range data.Items {
			if strings.Contains(item.Name, req.Keyword) {
				configMapList = append(configMapList, item)
			}
		}
	} else {
		configMapList = data.Items
	}

	result, err := paginate.Paginate(configMapList, req.Page, req.PageSize)

	return result.(*[]v1.ConfigMap), len(configMapList), nil

}
func (k *K8sConfigMapService) DescribeConfigMap(req configmap.DescribeConfigMapReq, uuid uuid.UUID) (*v1.ConfigMap, error) {
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
	configMapIns, err := client.CoreV1().ConfigMaps(req.Namespace).Get(context.TODO(), req.ConfigMapName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return configMapIns, nil
}
func (k *K8sConfigMapService) UpdateConfigMap(req configmap.UpdateConfigMapReq, uuid uuid.UUID) (*v1.ConfigMap, error) {
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
	var configMapIns *v1.ConfigMap
	err = json.Unmarshal([]byte(tmp), &configMapIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.CoreV1().ConfigMaps(req.Namespace).Update(context.TODO(), configMapIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sConfigMapService) DeleteConfigMap(req configmap.DeleteConfigMapReq, uuid uuid.UUID) error {
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

	err = client.CoreV1().ConfigMaps(req.Namespace).Delete(context.TODO(), req.ConfigMapName, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())

		return err
	}
	return nil
}
func (k *K8sConfigMapService) CreateConfigMap(req configmap.CreateConfigMapReq, uuid uuid.UUID) (*v1.ConfigMap, error) {
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
	var configmapIns *v1.ConfigMap
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &configmapIns)
	ins, err := client.CoreV1().ConfigMaps(req.Namespace).Create(context.TODO(), configmapIns, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
