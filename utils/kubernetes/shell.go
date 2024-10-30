package kubernetes

import (
	"DYCLOUD/global"
	cluster2 "DYCLOUD/model/kubernetes/cluster"
	"DYCLOUD/model/kubernetes/pods"
	"DYCLOUD/model/kubernetes/ws"
	request2 "DYCLOUD/model/system/request"
	"bytes"
	"encoding/json"
	"fmt"
	uuid "github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
	"sync"
	"time"
)

type ShellMessage struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
	Rows uint16 `json:"rows,omitempty"`
	Cols uint16 `json:"cols,omitempty"`
}

type RecordData struct {
	Event string  `json:"event"` // 输入输出事件
	Time  float64 `json:"time"`  // 时间差
	Data  []byte  `json:"data"`  // 数据
}

type KubeShell struct {
	sync.RWMutex
	Conn      *websocket.Conn
	sizeChan  chan remotecommand.TerminalSize
	stopChan  chan struct{}
	tty       bool
	Recorder  []*RecordData // 操作记录
	CreatedAt time.Time     // 创建时间
	UpdatedAt time.Time     // 最新的更新时间
	WriteAt   time.Time     // 写入时间
	written   bool          // 是否已写入记录, 一个流只允许写入一次
}

var EOT = "\u0004"
var _ TtyHandler = &KubeShell{}

func NewKubeShell(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*KubeShell, error) {
	// 升级get请求为websocket协议
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}
	kubeShell := &KubeShell{
		Conn:      conn,
		sizeChan:  make(chan remotecommand.TerminalSize),
		stopChan:  make(chan struct{}),
		tty:       true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Recorder:  make([]*RecordData, 0),
	}

	return kubeShell, nil
}

func (k *KubeShell) Stdin() io.Reader {
	return k
}

func (k *KubeShell) Stdout() io.Writer {
	return k
}

func (k *KubeShell) Stderr() io.Writer {
	return k
}

func (k *KubeShell) Tty() bool {
	return k.tty
}

func (k *KubeShell) Next() *remotecommand.TerminalSize {
	select {
	case size := <-k.sizeChan:
		return &size
	case <-k.stopChan:
		return nil
	}
}

func (k *KubeShell) Done() {
	close(k.stopChan)
}

func (k *KubeShell) Close() error {
	return k.Conn.Close()
}

func (k *KubeShell) Read(p []byte) (n int, err error) {
	k.WriteAt = time.Now()
	k.UpdatedAt = time.Now() // 更新时间
	_, message, err := k.Conn.ReadMessage()
	if err != nil {
		return copy(p, EOT), err
	}

	var msg ShellMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		return copy(p, EOT), err
	}

	switch msg.Type {
	case "read":
		return copy(p, msg.Data), nil
	case "resize":
		k.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	default:
		return copy(p, EOT), fmt.Errorf("unknown message type: %s", msg.Type)
	}

}

func (k *KubeShell) Write(p []byte) (n int, err error) {
	k.UpdatedAt = time.Now() // 更新时间
	msg, err := json.Marshal(ShellMessage{
		Type: "write",
		Data: string(p),
	})
	if err != nil {
		return 0, err
	}

	var data = make([]byte, len(p))
	copy(data, p)
	k.Recorder = append(k.Recorder, &RecordData{
		Time:  time.Since(k.CreatedAt).Seconds(),
		Event: "o",
		Data:  data,
	})
	if err := k.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (k *KubeShell) WriteLog(terminal ws.TerminalRequest, user *request2.CustomClaims) (err error) {
	// 记录用户的操作
	k.Lock()

	defer k.Unlock()

	if k.written {
		return nil
	}

	recorders := k.Recorder
	if len(recorders) != 0 {
		b := new(bytes.Buffer)
		meta := CastV2Header{
			Width:     160,
			Height:    24,
			Timestamp: time.Now().Unix(),
			Title:     "",
			Env: &map[string]string{
				"SHELL": "/bin/bash", "TERM": "xterm-256color",
			},
		}

		cast, buffer := NewCastV2(meta, b)
		for _, v := range recorders {
			cast.Record(v.Time, v.Data, v.Event)
		}

		var cluster cluster2.K8sCluster
		if err := global.DYCLOUD_DB.Where("id = ?", terminal.ClusterId).First(&cluster).Error; err != nil {
			global.DYCLOUD_LOG.Error("cluster get failed, err:", zap.Any("err", err))
			return err
		}

		terminalAudit := pods.PodRecord{
			Cluster:       cluster.Name,
			Namespace:     terminal.Namespace,
			PodName:       terminal.PodName,
			ContainerName: terminal.Name,
			Username:      user.Username,
			NickName:      user.NickName,
			UUID:          uuid.Must(uuid.NewV4()),
			Records:       DoZlibCompress(buffer.Bytes()),
		}

		if err = global.DYCLOUD_DB.Save(&terminalAudit).Error; err != nil {
			global.DYCLOUD_LOG.Error("terminalAudit save failed, err:", zap.Any("err", err))

		}
	}

	return err
}
