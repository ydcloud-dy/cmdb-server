package cloudtty

import (
	"DYCLOUD/global"
	model "DYCLOUD/model/kubernetes/cluster"
	kubernetesReq "DYCLOUD/model/kubernetes/cluster/request"
	cluster2 "DYCLOUD/service/kubernetes/cluster"
	"DYCLOUD/utils/kubernetes"
	"bytes"
	"encoding/json"
	"fmt"
	uuid "github.com/gofrs/uuid/v5"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/url"
)

type NodeTTYService struct{}

// @function: NodeshellPod
// @description: 获取nodeshell pod 信息
// @param: user model.User, cluster model.Cluster, nodetty kubernetesReq.NodeTTY
// @return: nodeshell []byte, podName string, err error
func (t *NodeTTYService) NodeshellPod(user model.User, cluster model.K8sCluster, nodetty kubernetesReq.NodeTTY) (nodeshell []byte, podName string, err error) {
	podName = fmt.Sprintf("node-shell-%s-%s", user.Username, randomString())
	// Define the pod object
	podData := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: "default",
			Labels: map[string]string{
				"devops/cluster-id": cluster.UUID.String(),
				"devops/manage":     "devops",
				"devops/username":   user.Username,
				"devops/podName":    podName,
			},
		},
		Spec: corev1.PodSpec{
			ActiveDeadlineSeconds: func(i int64) *int64 { return &i }(1800),
			Containers: []corev1.Container{
				{
					Name:  "node-shell",
					Image: "swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/alpine:latest",
					Command: []string{
						"nsenter",
					},
					Args: []string{
						"-t", "1",
						"-m",
						"-u",
						"-i",
						"-n",
						"sleep",
						"14000",
					},
					TerminationMessagePath:   "/dev/termination-log",
					TerminationMessagePolicy: corev1.TerminationMessagePolicy("File"),
					ImagePullPolicy:          corev1.PullIfNotPresent,
					SecurityContext: &corev1.SecurityContext{
						Privileged: new(bool),
					},
				},
			},
			RestartPolicy:                 corev1.RestartPolicyNever,
			TerminationGracePeriodSeconds: new(int64),
			DNSPolicy:                     corev1.DNSClusterFirst,
			NodeName:                      nodetty.NodeName,
			HostNetwork:                   true,
			HostPID:                       true,
			HostIPC:                       true,
			SecurityContext:               &corev1.PodSecurityContext{},
			SchedulerName:                 "default-scheduler",
			Tolerations: []corev1.Toleration{
				{
					Operator: corev1.TolerationOpExists,
				},
			},
			PriorityClassName:  "system-node-critical",
			Priority:           new(int32),
			EnableServiceLinks: new(bool),
			PreemptionPolicy:   new(corev1.PreemptionPolicy),
		},
	}

	*podData.Spec.Containers[0].SecurityContext.Privileged = true
	*podData.Spec.TerminationGracePeriodSeconds = 0
	*podData.Spec.Priority = 2000001000
	*podData.Spec.EnableServiceLinks = true
	preemptionPolicy := corev1.PreemptLowerPriority
	*podData.Spec.PreemptionPolicy = preemptionPolicy

	body, err := json.Marshal(podData)
	if err != nil {
		global.DYCLOUD_LOG.Error("Error JSON Marshal: " + err.Error())
		return body, podName, err
	}

	return body, podName, err
}

//@function: CreateNodeshell
//@description: 创建 pod 信息
//@param: user model.User, cluster model.Cluster, nodetty kubernetesReq.NodeTTY
//@return: nodeshell []byte, podName string, err error

func (t *NodeTTYService) CreateNodeshell(cluster model.K8sCluster, body []byte, httpClient *http.Client) (err error) {
	//url格式化
	apiUrl, err := url.Parse(fmt.Sprintf("%s%s", cluster.ApiAddress, "/api/v1/namespaces/default/pods"))
	if err != nil {
		global.DYCLOUD_LOG.Error("url Parse failed: " + err.Error())
		return err
	}

	req, err := http.NewRequest("POST", apiUrl.String(), bytes.NewBuffer(body))
	if err != nil {
		global.DYCLOUD_LOG.Error("http request failed: " + err.Error())
		return err
	}

	//发起请求
	_, err = httpClient.Do(req)
	if err != nil {
		global.DYCLOUD_LOG.Error("http request do failed: " + err.Error())
		return err
	}

	//取出数据Debug调试开启
	//rawResp, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(rawResp))

	return err
}

//@function: NodeTTY
//@description: 获取getNodeTTY
//@param:
//@return:

func (t *NodeTTYService) Get(nodetty kubernetesReq.NodeTTY, uuid uuid.UUID) (podMsg *PodMessage, err error) {
	clusterService := cluster2.K8sClusterService{}
	// 获取集群信息
	cluster, err := clusterService.GetK8sCluster(nodetty.ClusterId)
	if err != nil {
		fmt.Println(err.Error())
		global.DYCLOUD_LOG.Error("get cluster info failed: " + err.Error())
		return podMsg, err
	}

	// 获取集群用户信息
	user, err := clusterService.GetClusterByUserUUID(nodetty.ClusterId, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("get user info failed: " + err.Error())
		return podMsg, err
	}

	// 生成transport
	ts, err := kubernetes.GenerateTLSTransport(&cluster, &user, false)
	if err != nil {
		global.DYCLOUD_LOG.Error("transport generate failed: " + err.Error())
		return podMsg, err
	}

	// 生成httpClient
	httpClient := &http.Client{Transport: ts}

	// 获取Node Shell Pod信息
	nodeshell, pod_name, err := t.NodeshellPod(user, cluster, nodetty)
	if err != nil {
		global.DYCLOUD_LOG.Error("Node shell Pod Message failed: " + err.Error())
		return podMsg, err
	}

	// 创建 CreateNodeshell
	if err = t.CreateNodeshell(cluster, nodeshell, httpClient); err != nil {
		global.DYCLOUD_LOG.Error("Create Cloudshell failed: " + err.Error())
		return podMsg, err
	}

	return &PodMessage{
		Name:      pod_name,
		Namespace: "default",
		Container: "node-shell",
	}, err
}
