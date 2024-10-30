package service

import (
	global2 "DYCLOUD/global/docker"
	model "DYCLOUD/model/docker"
	"DYCLOUD/utils/docker/docker"
	"context"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"sort"
	"strings"
)

type VolumeService struct {
}

// ListVolume 获取存储卷
func (e *VolumeService) ListVolume(host string, req model.SearchVolume) (res []model.Volume, err error) {

	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return nil, err
	}

	option := volume.ListOptions{}

	if strings.Trim(req.Name, "") != "" {
		option.Filters = filters.NewArgs(filters.KeyValuePair{Key: "name", Value: req.Name})
	}

	var volumeResponse volume.ListResponse

	volumeResponse, err = cli.VolumeList(context.TODO(), option)
	if err != nil {
		return nil, err
	}

	volumeList := make([]model.Volume, 0)
	for _, re := range volumeResponse.Volumes {

		volumeList = append(volumeList, model.Volume{
			Name:       re.Name,
			Driver:     re.Driver,
			Created:    docker.StringTimeFormat(re.CreatedAt),
			MountPoint: re.Mountpoint,
			Status:     re.Status,
			Labels:     re.Labels,
			Options:    re.Options,
		})
	}

	sort.Slice(volumeList, func(i, j int) bool {
		if docker.StringToTimestamp(volumeList[i].Created) > docker.StringToTimestamp(volumeList[j].Created) {
			return true
		} else {
			return false
		}
	})

	return volumeList, nil
}

// RemoveVolume 删除存储卷
func (e *VolumeService) RemoveVolume(host string, ids []string) error {
	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return err
	}

	for _, id := range ids {
		err := cli.VolumeRemove(context.Background(), id, true)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateVolume 创建存储卷
func (e *VolumeService) CreateVolume(host string, n model.Volume) error {

	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return err
	}

	option := volume.CreateOptions{
		Name:   n.Name,
		Driver: n.Driver,
	}

	if n.Cifs.Enable {
		option.DriverOpts = map[string]string{
			"device": n.Cifs.Device,
			"o":      n.Cifs.Option(),
			"type":   "cifs",
		}
	} else if n.Nfs.Enable {
		option.DriverOpts = map[string]string{
			"device": n.Nfs.Device,
			"o":      n.Nfs.Option(),
			"type":   "nfs",
		}
	}

	_, err = cli.VolumeCreate(context.Background(), option)
	if err != nil {
		return err
	}

	return nil
}
