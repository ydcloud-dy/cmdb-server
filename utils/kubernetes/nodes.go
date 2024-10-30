package kubernetes

import (
	"DYCLOUD/global"
	cluster2 "DYCLOUD/model/kubernetes/cluster"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

type NodeService interface {
	MetricsConfig() (config *rest.Config, err error)
	Client() (*metricsv.Clientset, error)
}

// Kubernetes
// @Description: 结构体包含 k8sCluster集群对象、集群用户信息和是否为管理员的标志
type NodeMetrics struct {
	*cluster2.K8sCluster
	*cluster2.User
	IsAdmin bool
}

// NewKubernetes
//
//	@Description: 创建一个新的NodeMetrics实例
//	@return *Kubernetes
func NewNodeMetrics(cluster *cluster2.K8sCluster, user *cluster2.User, isAdmin bool) *NodeMetrics {
	return &NodeMetrics{cluster, user, isAdmin}

}
func (k *NodeMetrics) Client() (*metricsv.Clientset, error) {
	cfg, err := k.MetricsConfig()
	if err != nil {
		global.DYCLOUD_LOG.Error("Config init failed: ", zap.Error(err))
		return nil, err
	}
	clientset, err := metricsv.NewForConfig(cfg)
	if err != nil {
		global.DYCLOUD_LOG.Error("clientset init failed: ", zap.Error(err))
		return nil, err
	}

	return clientset, err
}

// Config
//
//	@Description: 初始化k8s配置
//	@receiver k
//	@return config
//	@return err
func (k *NodeMetrics) MetricsConfig() (config *rest.Config, err error) {
	// 判断用户创建k8s集群的方式，1为kubeconfig，2为token
	if k.KubeType == 1 {
		if k.IsAdmin {
			config, err = clientcmd.RESTConfigFromKubeConfig([]byte(k.K8sCluster.KubeConfig))
			if err != nil {
				global.DYCLOUD_LOG.Error("RESTConfigFromKubeConfig init failed: ", zap.Error(err))
				return nil, err
			}
		} else {
			config, err = clientcmd.RESTConfigFromKubeConfig([]byte(k.User.KubeConfig))
			if err != nil {
				global.DYCLOUD_LOG.Error("RESTConfigFromKubeConfig init failed: ", zap.Error(err))
				return nil, err
			}
		}
		return config, err
	}
	if k.KubeType == 2 {
		if k.IsAdmin {
			return &rest.Config{
				Host:            k.ApiAddress,
				BearerToken:     k.K8sCluster.KubeConfig,
				TLSClientConfig: rest.TLSClientConfig{Insecure: true},
			}, err
		} else {
			return &rest.Config{
				Host:            k.ApiAddress,
				BearerToken:     k.User.KubeConfig,
				TLSClientConfig: rest.TLSClientConfig{Insecure: true},
			}, err
		}
	}

	return
}
