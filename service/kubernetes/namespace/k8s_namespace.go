package namespace

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/namespaces"
	"DYCLOUD/service/kubernetes/cluster"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sNamespaceService struct {
	kubernetes.BaseService
}

func (k *K8sNamespaceService) GetNamespaceList(req namespaces.GetNamespaceListReq, uuid uuid.UUID) (*[]v1.Namespace, int, error) {
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
	var s = cluster.K8sClusterService{}

	nsList, err := s.GetClusterUserNamespace(req.ClusterId, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	options := metav1.ListOptions{}
	data, err := client.CoreV1().Namespaces().List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	// 获取用户有的命名空间列表
	var userNS []v1.Namespace
	for _, nsName := range nsList {
		for _, nsDetail := range data.Items {
			if nsDetail.Name == nsName {
				userNS = append(userNS, nsDetail)
			}
		}
	}
	var filterNamespace []v1.Namespace
	if req.Keyword != "" {
		for _, item := range userNS {
			if strings.Contains(item.Name, req.Keyword) {
				filterNamespace = append(filterNamespace, item)
			}
		}
	} else {
		filterNamespace = userNS
	}

	result, err := paginate.Paginate(filterNamespace, req.Page, req.PageSize)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}

	return result.(*[]v1.Namespace), len(filterNamespace), nil
}
func (k *K8sNamespaceService) DescribeNamespace(req namespaces.DescribeNamespaceReq, uuid uuid.UUID) (*v1.Namespace, error) {
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
	namespaceIns, err := client.CoreV1().Namespaces().Get(context.TODO(), req.NamespaceName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return namespaceIns, nil
}
func (k *K8sNamespaceService) UpdateNamespace(req namespaces.UpdateNamespaceReq, uuid uuid.UUID) (*v1.Namespace, error) {
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
	var namespaceIns *v1.Namespace
	err = json.Unmarshal([]byte(tmp), &namespaceIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.CoreV1().Namespaces().Update(context.TODO(), namespaceIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sNamespaceService) DeleteNamespace(req namespaces.DeleteNamespaceReq, uuid uuid.UUID) error {
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

	err = client.CoreV1().Namespaces().Delete(context.TODO(), req.NamespaceName, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())

		return err
	}
	return nil
}
func (k *K8sNamespaceService) CreateNamespace(req namespaces.CreateNamespaceReq, uuid uuid.UUID) (*v1.Namespace, error) {
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
	var namespace *v1.Namespace
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &namespace)
	ins, err := client.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
