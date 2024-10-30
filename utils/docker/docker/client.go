package docker

import (
	"DYCLOUD/global"
	"DYCLOUD/model/docker"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/docker/docker/client"
	"net/http"
	"os"
	"sync"
)

type DockerClient struct {
	storage sync.Map
}

// NewDockerClient 创建Docker客户端
func NewDockerClient() (d *DockerClient) {
	d = &DockerClient{
		storage: sync.Map{},
	}
	return
}

func (d *DockerClient) Load(host string) (cli *client.Client, err error) {
	v, ok := d.storage.Load(host)
	if ok {
		cli = v.(*client.Client)
	} else {
		var hs model.Host
		err = global.DYCLOUD_DB.Where("name=?", host).First(&hs).Error
		if err != nil {
			global.DYCLOUD_LOG.Error(err.Error())
			return nil, errors.New("host is not found")
		}
		cli, err = d.CreateClient(hs)
		if err != nil {
			global.DYCLOUD_LOG.Error(err.Error())
			return nil, err
		}
		d.storage.Store(host, cli)
	}
	return
}

func (d *DockerClient) CreateClient(hs model.Host) (cli *client.Client, err error) {
	switch hs.Type {
	case model.HostTypeApi:
		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: *hs.SkipCert,
				},
			},
		}
		if *hs.EnableTls {
			caPath, err2 := d.SaveCertFile(fmt.Sprintf("%v-%v", hs.Name, "ca"), hs.TlsCa)
			if err2 != nil {
				return nil, err2
			}
			certPath, err := d.SaveCertFile(fmt.Sprintf("%v-%v", hs.Name, "cert"), hs.TlsCert)
			if err != nil {
				return nil, err
			}
			keyPath, err := d.SaveCertFile(fmt.Sprintf("%v-%v", hs.Name, "key"), hs.TlsKey)
			if err != nil {
				return nil, err
			}
			cli, err = client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.WithHTTPClient(httpClient), client.WithHost(fmt.Sprintf("tcp://%v:%v", hs.ApiAddress, *hs.Port)), client.WithTLSClientConfig(caPath, certPath, keyPath))
		} else {
			fmt.Println(fmt.Sprintf("tcp://%v:%v", hs.ApiAddress, *hs.Port))
			cli, err = client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.WithHost(fmt.Sprintf("tcp://%v:%v", hs.ApiAddress, *hs.Port)))
		}
	case model.HostTypeSocket:

		cli, err = client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.WithHost(fmt.Sprintf("unix://%v", hs.SocketPath)))
	}
	return cli, err
}

func (d *DockerClient) SaveCertFile(key string, context string) (path string, err error) {
	filename := fmt.Sprintf("%v/%v", os.TempDir(), key)
	err = os.WriteFile(filename, []byte(context), 0644)
	return filename, err
}

func (d *DockerClient) Remove(host string) {
	if _, ok := d.storage.Load(host); ok {
		d.storage.Delete(host)
	}
}
