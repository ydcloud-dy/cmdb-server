package service

import (
	global2 "DYCLOUD/global/docker"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var SystemServiceApp = new(SystemService)

type SystemService struct {
	dockerClient *client.Client
}

// Info 获取系统信息
func (s SystemService) Info(host string) (types.Info, error) {

	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return types.Info{}, err
	}

	info, err := cli.Info(context.TODO())
	if err != nil {
		return types.Info{}, err
	}
	return info, nil
}
