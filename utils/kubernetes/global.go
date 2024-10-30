package kubernetes

import (
	"DYCLOUD/global"
	model "DYCLOUD/model/kubernetes/cluster"
	"bytes"
	"compress/zlib"
	"go.uber.org/zap"
	"io"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"unsafe"
)

type kubeclient struct {
	Pod *Pod
}

func NewKubeClient(id int) (kube *kubeclient, err error) {
	var cluster model.K8sCluster
	if err := global.DYCLOUD_DB.Where("id = ?", id).First(&cluster).Error; err != nil {
		global.DYCLOUD_LOG.Error("cluster get failed, err:", zap.Any("err", err))
		return nil, err
	}

	var config *rest.Config
	if cluster.KubeType == 1 {
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.KubeConfig))
		if err != nil {
			global.DYCLOUD_LOG.Error("config get failed, err:", zap.Any("err", err))
		}
	} else if cluster.KubeType == 2 {
		config = &rest.Config{
			Host:            cluster.ApiAddress,
			BearerToken:     cluster.KubeConfig,
			TLSClientConfig: rest.TLSClientConfig{Insecure: true},
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		global.DYCLOUD_LOG.Error("kubernetes client init failed, err:", zap.Any("err", err))
	}

	client := &kubeclient{
		Pod: NewPod(clientset, config),
	}

	return client, nil
}

func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ZlibCompress 进行zlib压缩
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

// ZlibCompress 进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

func Str2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
