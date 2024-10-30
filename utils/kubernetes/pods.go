package kubernetes

import (
	"DYCLOUD/model/kubernetes/ws"
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type Pod struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
}

func NewPod(clientset *kubernetes.Clientset, config *rest.Config) *Pod {
	return &Pod{
		clientset: clientset,
		config:    config,
	}
}

func (cli *Pod) ContainerLog(kubeLogger *KubeLogger, name, pod_name, namespace string) {
	lines := int64(500)
	opts := corev1.PodLogOptions{Container: name, Follow: true, TailLines: &lines}

	err := cli.LogsStream(pod_name, namespace, &opts, kubeLogger)
	if err != nil {
		kubeLogger.Write([]byte(err.Error()))
	}
}

func (cli *Pod) Logs(name, namespace string, ops *corev1.PodLogOptions) *rest.Request {
	return cli.clientset.CoreV1().Pods(namespace).GetLogs(name, ops)
}

func (cli *Pod) LogsStream(name, namespace string, opts *corev1.PodLogOptions, write io.Writer) error {
	req := cli.Logs(name, namespace, opts)
	stream, err := req.Stream(context.TODO())
	if err != nil {
		return err
	}

	defer stream.Close()

	buf := bufio.NewReader(stream)
	for {
		//一直从buffer中读取数据
		bytes, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				_, err = write.Write(bytes)
			}
			return err

		}
		_, err = write.Write(bytes)
		if err != nil {
			_, err = write.Write(bytes)
		}
	}
}

type TtyHandler interface {
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
	Tty() bool
	remotecommand.TerminalSizeQueue
	Done()
}

func (cli *Pod) Exec(cmd []string, handler TtyHandler, t ws.TerminalRequest) error {
	defer func() {
		handler.Done()
	}()

	//构造请求
	req := cli.clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(t.Namespace).
		Name(t.PodName).SubResource("exec")
	req.VersionedParams(&corev1.PodExecOptions{
		Container: t.Name,
		Command:   cmd,
		Stdin:     handler.Stdin() != nil,
		Stdout:    handler.Stdout() != nil,
		Stderr:    handler.Stderr() != nil,
		TTY:       true,
	}, scheme.ParameterCodec)

	execustor, err := remotecommand.NewSPDYExecutor(cli.config, "POST", req.URL())
	if err != nil {
		return err
	}

	if err := execustor.Stream(remotecommand.StreamOptions{
		Stdin:             handler.Stdin(),
		Stdout:            handler.Stdout(),
		Stderr:            handler.Stdout(),
		Tty:               handler.Tty(),
		TerminalSizeQueue: handler,
	}); err != nil {
		return err
	}

	return err
}

func (cli *Pod) GetLogs(namespace, name string) (logs string, err error) {

	line := int64(500)
	ret := cli.clientset.CoreV1().Pods(namespace).GetLogs(name, &corev1.PodLogOptions{
		TypeMeta:  metav1.TypeMeta{},
		TailLines: &line,
	})
	podLogs, err := ret.Stream(context.TODO())
	if err != nil {
		return logs, errors.New("error in opening stream")
	}

	defer podLogs.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return logs, errors.New("error in copy information from podLogs to buf")
	}

	str := buf.String()
	return str, err
}
