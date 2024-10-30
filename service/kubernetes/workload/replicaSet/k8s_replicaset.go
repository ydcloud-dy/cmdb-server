package replicaSet

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/replicaSet"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sReplicaSetService struct {
	kubernetes.BaseService
}

func (k *K8sReplicaSetService) GetReplicaSetList(req replicaSet.GetReplicaSetListReq, uuid uuid.UUID) (*[]v1.ReplicaSet, int, error) {
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
	data, err := client.AppsV1().ReplicaSets(req.Namespace).List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterReplicaSet []v1.ReplicaSet
	if req.Keyword != "" {
		for _, item := range data.Items {
			if strings.Contains(item.Name, req.Keyword) {
				filterReplicaSet = append(filterReplicaSet, item)
			}
		}
	} else {
		filterReplicaSet = data.Items
	}
	if req.Page == 0 || req.PageSize == 0 {
		return &filterReplicaSet, len(filterReplicaSet), nil
	}
	result, err := paginate.Paginate(filterReplicaSet, req.Page, req.PageSize)

	return result.(*[]v1.ReplicaSet), len(filterReplicaSet), nil
}
