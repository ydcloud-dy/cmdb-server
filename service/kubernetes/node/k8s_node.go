package node

import (
	"DYCLOUD/global"
	cluster2 "DYCLOUD/model/kubernetes/cluster"
	"DYCLOUD/model/kubernetes/nodes"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sClientSet "k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/drain"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"os"
	"strings"
	"time"
)

type K8sNodeService struct {
	kubernetes.BaseService
}

func GetClusterByUserUUID(id int, uuid uuid.UUID) (user cluster2.User, err error) {
	err = global.DYCLOUD_DB.Where("cluster_id = ? and uuid = ?", id, uuid).First(&user).Error
	return
}

// GetNodeList
//
// @Description: 获取node列表
// @receiver k
// @param req query nodes.NodeListReq true "节点列表请求参数"
// @param uuid path uuid.UUID true "用户UUID"
// @return *corev1.NodeList
// @return error
func (k *K8sNodeService) GetNodeList(req nodes.NodeListReq, uuid uuid.UUID) (node *[]corev1.Node, total int, err error) {

	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &[]corev1.Node{}, 0, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &[]corev1.Node{}, 0, err
	}
	options := metav1.ListOptions{}
	nodeList, err := client.CoreV1().Nodes().List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &[]corev1.Node{}, 0, err
	}
	var filterNode []corev1.Node
	if req.Keyword != "" {
		for _, item := range nodeList.Items {
			if strings.Contains(item.Name, req.Keyword) {
				filterNode = append(filterNode, item)
			}
		}
	} else {
		filterNode = nodeList.Items
	}

	result, err := paginate.Paginate(filterNode, req.Page, req.PageSize)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &[]corev1.Node{}, 0, err
	}

	return result.(*[]corev1.Node), len(nodeList.Items), nil

}

// GetNodeMetricsList
//
// @Description: 获取node的metrics列表
// @receiver k
// @param req query nodes.NodeMetricsReq true "节点Metrics请求参数"
// @param uuid path uuid.UUID true "用户UUID"
// @return *v1beta1.NodeMetricsList
// @return error
func (k *K8sNodeService) GetNodeMetricsList(req nodes.NodeMetricsReq, uuid uuid.UUID) (*v1beta1.NodeMetricsList, error) {
	_, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &v1beta1.NodeMetricsList{}, err
	}
	// 根据集群实例生成MetricsClientset对象
	client, err := k.NodeMetrics.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &v1beta1.NodeMetricsList{}, err
	}
	// 查询node Metrics列表并返回
	nodeMetrics, err := client.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &v1beta1.NodeMetricsList{}, nil
	}
	return nodeMetrics, nil
}

// DescribeNodeInfo
//
// @Description: 获取节点详情信息
// @receiver k
// @param req query nodes.DescribeNodeReq true "节点详情请求参数"
// @param uuid path uuid.UUID true "用户UUID"
// @return *corev1.Node
// @return error
func (k *K8sNodeService) DescribeNodeInfo(req nodes.DescribeNodeReq, uuid uuid.UUID) (*corev1.Node, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &corev1.Node{}, err
	}
	// 根据集群实例生成Clientset对象
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &corev1.Node{}, err
	}
	// 查询node 详情并返回
	nodeInfo, err := client.CoreV1().Nodes().Get(context.TODO(), req.NodeName, metav1.GetOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &corev1.Node{}, err
	}
	nodeInfo.APIVersion = "v1"
	nodeInfo.Kind = "Node"
	return nodeInfo, nil
}

// UpdateNodeInfo
//
// @Description: 更新节点信息
// @receiver k
// @param req body nodes.UpdateNodeReq true "更新节点请求参数"
// @param uuid path uuid.UUID true "用户UUID"
// @return *corev1.Node
// @return error
func (k *K8sNodeService) UpdateNodeInfo(req nodes.UpdateNodeReq, uuid uuid.UUID) (*corev1.Node, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &corev1.Node{}, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &corev1.Node{}, err
	}

	var nodeIns = &corev1.Node{}
	s, err := json.Marshal(&req.Content)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &corev1.Node{}, err
	}
	json.Unmarshal(s, &nodeIns)
	//err = json.Unmarshal([]byte(req.Content), &nodeIns)
	//if err != nil {
	//	global.DYCLOUD_LOG.Error(err.Error())
	//	return &corev1.Node{}, err
	//}
	nodeInfo, err := client.CoreV1().Nodes().Update(context.TODO(), nodeIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &corev1.Node{}, err
	}
	return nodeInfo, nil
}

// EvictAllNodePod
//
// @Description: 驱逐节点上的所有 Pod
// @receiver k
// @param req body nodes.EvictAllNodePodReq true "驱逐节点上所有Pod的请求参数"
// @return error
func (k *K8sNodeService) EvictAllNodePod(req nodes.EvictAllNodePodReq) error {
	// 先到库里根据clusterID查询到集群实例
	var clusterIns = &cluster2.K8sCluster{}
	if err := global.DYCLOUD_DB.Where("id = ?", req.ClusterId).First(clusterIns).Error; err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}
	// 根据集群实例生成Clientset对象
	kubernetes := kubernetes.NewKubernetes(clusterIns, &cluster2.User{}, true)
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}
	node, err := client.CoreV1().Nodes().Get(context.TODO(), req.NodeName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	node.Spec.Unschedulable = true
	_, err = client.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})

	err = drainNode(client, req.NodeName)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}
	return nil

}

// drainNode 驱逐节点上的所有 Pod
func drainNode(client *k8sClientSet.Clientset, nodeName string) error {
	//node, err := client.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	//if err != nil {
	//	return err
	//}

	drainer := &drain.Helper{
		Ctx:                 context.TODO(),
		Client:              client,
		Force:               true,
		IgnoreAllDaemonSets: true,
		DeleteEmptyDirData:  true,
		Timeout:             5 * time.Minute,
		Out:                 os.Stdout,
		ErrOut:              os.Stderr,
	}

	// 这里不再调用 cordon，因为在 drain 过程中会自动进行
	if err := drain.RunNodeDrain(drainer, nodeName); err != nil {
		return err
	}

	return nil
}

// evictPod 驱逐指定的 Pod
func evictPod(client *k8sClientSet.Clientset, pod *corev1.Pod) error {
	deleteOptions := &metav1.DeleteOptions{}
	eviction := &policyv1.Eviction{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
		DeleteOptions: deleteOptions,
	}
	err := client.PolicyV1().Evictions(eviction.Namespace).Evict(context.TODO(), eviction)
	if err != nil {
		return err
	}
	//var wg sync.WaitGroup
	//wg.Add(1)
	//defer wg.Done()
	//// 等待 Pod 驱逐完成
	//go func() {
	//	for {
	//		_, err := client.CoreV1().Pods(pod.Namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
	//		if err != nil {
	//			break
	//		}
	//		time.Sleep(1 * time.Second)
	//	}
	//}()
	//wg.Wait()

	return nil
}
