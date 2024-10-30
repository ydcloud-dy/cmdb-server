package kubernetes

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type KubeLogger struct {
	Conn *websocket.Conn
}

func NewKubeLogger(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*KubeLogger, error) {
	//升级get 请求
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}

	return &KubeLogger{Conn: conn}, nil
}

func (kl *KubeLogger) Write(data []byte) (n int, err error) {
	if err := kl.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return 0, err
	}
	return len(data), nil
}

func (kl *KubeLogger) Close() error {
	return kl.Close()
}
