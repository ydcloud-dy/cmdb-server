package horizontalPod

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/horizontalPod"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sHorizontalPodService struct {
	kubernetes.BaseService
}

func (k *K8sHorizontalPodService) GetHorizontalPodList(req horizontalPod.GetHorizontalPodListReq, uuid uuid.UUID) (*[]v1.HorizontalPodAutoscaler, int, error) {
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
	data, err := client.AutoscalingV1().HorizontalPodAutoscalers(req.Namespace).List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterHorizontalPod []v1.HorizontalPodAutoscaler
	if req.Keyword != "" {
		for _, item := range data.Items {
			if strings.Contains(item.Name, req.Keyword) {
				filterHorizontalPod = append(filterHorizontalPod, item)
			}
		}
	} else {
		filterHorizontalPod = data.Items
	}

	result, err := paginate.Paginate(filterHorizontalPod, req.Page, req.PageSize)

	return result.(*[]v1.HorizontalPodAutoscaler), len(filterHorizontalPod), nil
}
func (k *K8sHorizontalPodService) DescribeHorizontalPod(req horizontalPod.DescribeHorizontalPodReq, uuid uuid.UUID) (*v1.HorizontalPodAutoscaler, error) {
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
	HorizontalPodIns, err := client.AutoscalingV1().HorizontalPodAutoscalers(req.Namespace).Get(context.TODO(), req.HorizontalPodName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return HorizontalPodIns, nil
}
func (k *K8sHorizontalPodService) UpdateHorizontalPod(req horizontalPod.UpdateHorizontalPodReq, uuid uuid.UUID) (*v1.HorizontalPodAutoscaler, error) {
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
	var HorizontalPodIns *v1.HorizontalPodAutoscaler
	err = json.Unmarshal([]byte(tmp), &HorizontalPodIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.AutoscalingV1().HorizontalPodAutoscalers(req.Namespace).Update(context.TODO(), HorizontalPodIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sHorizontalPodService) DeleteHorizontalPod(req horizontalPod.DeleteHorizontalPodReq, uuid uuid.UUID) error {
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

	err = client.AutoscalingV1().HorizontalPodAutoscalers(req.Namespace).Delete(context.TODO(), req.HorizontalPodName, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())

		return err
	}
	return nil
}
func (k *K8sHorizontalPodService) CreateHorizontalPod(req horizontalPod.CreateHorizontalPodReq, uuid uuid.UUID) (*v1.HorizontalPodAutoscaler, error) {
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
	var HorizontalPod *v1.HorizontalPodAutoscaler
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &HorizontalPod)
	ins, err := client.AutoscalingV1().HorizontalPodAutoscalers(req.Namespace).Create(context.TODO(), HorizontalPod, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
