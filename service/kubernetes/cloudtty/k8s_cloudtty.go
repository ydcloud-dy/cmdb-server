package cloudtty

import (
	"DYCLOUD/global"
	"DYCLOUD/model/cloudTTY"
	cluster3 "DYCLOUD/model/kubernetes/cluster"
	cluster2 "DYCLOUD/service/kubernetes/cluster"
	"DYCLOUD/utils/kubernetes"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/rand"
	"net/http"
	"net/url"
	"sigs.k8s.io/yaml"
	"time"
)

type K8sCloudTTYService struct{}
type CloudShell struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
}
type Metadata struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
}

type Spec struct {
	ConfigmapName string `json:"configmapName"`
	CommandAction string `json:"commandAction"`
	Cleanup       bool   `json:"cleanup"`
	Image         string `json:"image"`
}
type PodMessage struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Container string `json:"container"`
}

func (k *K8sCloudTTYService) Get(tty cloudTTY.CloudTTY, uuid uuid.UUID) (podMsg *PodMessage, err error) {
	c := cluster2.K8sClusterService{}
	cluster, err := c.GetK8sCluster(tty.ClusterId)
	if err != nil {
		global.DYCLOUD_LOG.Error("get cluster info failed: " + err.Error())
		return podMsg, err
	}
	// 获取集群信息
	user, err := c.GetClusterByUserUUID(tty.ClusterId, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("get cluster info failed: " + err.Error())
		return podMsg, err
	}

	// 获取集群用户信息
	// 获取ConfigMap 信息
	configmap, kubeconfigName, err := k.getConfigMap(user, cluster)
	if err != nil {
		global.DYCLOUD_LOG.Error("get ConfigMap Message failed: " + err.Error())
		return podMsg, err
	}

	// 获取 Cloudshell 信息
	cloudshell, cloudshellName, err := k.getCloudshell(user, cluster, kubeconfigName)
	if err != nil {
		global.DYCLOUD_LOG.Error("get cloudshell Message failed: " + err.Error())
		return podMsg, err
	}
	// 生成transport
	ts, err := kubernetes.GenerateTLSTransport(&cluster, &user, true)
	if err != nil {
		global.DYCLOUD_LOG.Error("transport generate failed: " + err.Error())
		return podMsg, err
	}

	// 生成httpClient
	httpClient := &http.Client{Transport: ts}

	// 删除ConfigMap
	if err = k.DeleteConfigMap(cluster, kubeconfigName, httpClient); err != nil {
		global.DYCLOUD_LOG.Error("Delete ConfigMap failed: " + err.Error())
		return podMsg, err
	}

	// 创建ConfigMap
	if err = k.CreateConfigMap(cluster, configmap, httpClient); err != nil {
		global.DYCLOUD_LOG.Error("Create ConfigMap failed: " + err.Error())
		return podMsg, err
	}

	// 创建 CreateCloudshell
	if err = k.CreateCloudshell(cluster, cloudshell, httpClient); err != nil {
		global.DYCLOUD_LOG.Error("Create Cloudshell failed: " + err.Error())
		return podMsg, err
	}

	return &PodMessage{Name: cloudshellName, Namespace: "default", Container: "web-tty"}, err
}
func (t *K8sCloudTTYService) CreateCloudshell(cluster cluster3.K8sCluster, body []byte, httpClient *http.Client) (err error) {
	//url格式化
	apiUrl, err := url.Parse(fmt.Sprintf("%s%s", cluster.ApiAddress, "/apis/cloudshell.cloudtty.io/v1alpha1/namespaces/default/cloudshells"))
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

func (t *K8sCloudTTYService) DeleteConfigMap(cluster cluster3.K8sCluster, configmap string, httpClient *http.Client) (err error) {
	//url格式化
	apiUrl, err := url.Parse(fmt.Sprintf("%s%s", cluster.ApiAddress, "/api/v1/namespaces/default/configmaps/"+configmap))
	if err != nil {
		global.DYCLOUD_LOG.Error("url Parse failed: " + err.Error())
		return err
	}

	req, err := http.NewRequest("DELETE", apiUrl.String(), nil)
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
func (t *K8sCloudTTYService) CreateConfigMap(cluster cluster3.K8sCluster, body []byte, httpClient *http.Client) (err error) {
	//url格式化
	apiUrl, err := url.Parse(fmt.Sprintf("%s%s", cluster.ApiAddress, "/api/v1/namespaces/default/configmaps"))
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

func (t *K8sCloudTTYService) getConfigMap(user cluster3.User, cluster cluster3.K8sCluster) (configmap []byte, kubeconfigName string, err error) {
	// 用户kubeconfig 凭据 JSON 转 Go 数据结构
	var jsonObj map[string]interface{}
	if err = json.Unmarshal([]byte(user.KubeConfig), &jsonObj); err != nil {
		global.DYCLOUD_LOG.Error("Error unmarshalling JSON: " + err.Error())
		return configmap, kubeconfigName, err
	}

	// 用户kubeconfig 数据结构转 YAML
	yamlData, err := yaml.Marshal(jsonObj)
	if err != nil {
		global.DYCLOUD_LOG.Error("Error marshalling to YAML: " + err.Error())
		return configmap, kubeconfigName, err
	}

	// 实例化configmap
	kubeconfigName = fmt.Sprintf("%s-kubeconfig-%s", user.Username, cluster.UUID)
	configMap := v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kubeconfigName,
			Namespace: "default",
			Labels: map[string]string{
				"devops/cluster-id": cluster.UUID.String(),
				"devops/manage":     "devops",
				"devops/username":   user.Username,
			},
		},
		Data: map[string]string{
			"config": string(yamlData),
		},
	}

	body, err := json.Marshal(configMap)
	if err != nil {
		global.DYCLOUD_LOG.Error("Error JSON Marshal: " + err.Error())
		return configmap, kubeconfigName, err
	}

	return body, kubeconfigName, err
}
func (t *K8sCloudTTYService) getCloudshell(user cluster3.User, cluster cluster3.K8sCluster, kubeconfigName string) (cloudshell []byte, podName string, err error) {
	podName = fmt.Sprintf("%s-%s", user.Username, randomString())
	cloudShellData := CloudShell{
		APIVersion: "cloudshell.cloudtty.io/v1alpha1",
		Kind:       "CloudShell",
		Metadata: Metadata{
			Name:      podName,
			Namespace: "default",
			Labels: map[string]string{
				"devops/cluster-id": cluster.UUID.String(),
				"devops/manage":     "devops",
				"devops/username":   user.Username,
				"devops/podName":    podName,
			},
		},
		Spec: Spec{
			ConfigmapName: kubeconfigName,
			CommandAction: "bash",
			Cleanup:       true,
			Image:         "ghcr.io/cloudtty/cloudshell:v0.6.0",
		},
	}

	body, err := json.Marshal(cloudShellData)
	if err != nil {
		global.DYCLOUD_LOG.Error("Error JSON Marshal: " + err.Error())
		return cloudshell, podName, err
	}

	return body, podName, err
}

// 随机生成字符串
func randomString() string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
