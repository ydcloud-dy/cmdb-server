package limitRange

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/limitRange"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sLimitRangeService struct {
	kubernetes.BaseService
}

func (k *K8sLimitRangeService) GetLimitRangeList(req limitRange.GetLimitRangeListReq, uuid uuid.UUID) (*[]v1.LimitRange, int, error) {
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
	data, err := client.CoreV1().LimitRanges(req.Namespace).List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterLimitRange []v1.LimitRange
	if req.Keyword != "" {
		for _, item := range data.Items {
			if strings.Contains(item.Name, req.Keyword) {
				filterLimitRange = append(filterLimitRange, item)
			}
		}
	} else {
		filterLimitRange = data.Items
	}

	result, err := paginate.Paginate(filterLimitRange, req.Page, req.PageSize)

	return result.(*[]v1.LimitRange), len(filterLimitRange), nil
}
func (k *K8sLimitRangeService) DescribeLimitRange(req limitRange.DescribeLimitRangeReq, uuid uuid.UUID) (*v1.LimitRange, error) {
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
	LimitRangeIns, err := client.CoreV1().LimitRanges(req.Namespace).Get(context.TODO(), req.LimitRangeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return LimitRangeIns, nil
}
func (k *K8sLimitRangeService) UpdateLimitRange(req limitRange.UpdateLimitRangeReq, uuid uuid.UUID) (*v1.LimitRange, error) {
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
	var LimitRangeIns *v1.LimitRange
	err = json.Unmarshal([]byte(tmp), &LimitRangeIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.CoreV1().LimitRanges(req.Namespace).Update(context.TODO(), LimitRangeIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sLimitRangeService) DeleteLimitRange(req limitRange.DeleteLimitRangeReq, uuid uuid.UUID) error {
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

	err = client.CoreV1().LimitRanges(req.Namespace).Delete(context.TODO(), req.LimitRangeName, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())

		return err
	}
	return nil
}
func (k *K8sLimitRangeService) CreateLimitRange(req limitRange.CreateLimitRangeReq, uuid uuid.UUID) (*v1.LimitRange, error) {
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
	var LimitRange *v1.LimitRange
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &LimitRange)
	ins, err := client.CoreV1().LimitRanges(req.Namespace).Create(context.TODO(), LimitRange, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
