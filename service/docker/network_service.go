package service

import (
	global2 "DYCLOUD/global/docker"
	model "DYCLOUD/model/docker"
	"DYCLOUD/utils/docker/docker"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"sort"
	"strings"
	"sync"
)

type NetworkService struct {
	store sync.Map
}

// ListNetwork 获取网络列表
func (e *NetworkService) ListNetwork(host string, req model.SearchNetwork) (res []model.Network, err error) {

	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return nil, err
	}

	option := types.NetworkListOptions{}

	if strings.Trim(req.Name, "") != "" {
		option.Filters = filters.NewArgs(filters.KeyValuePair{Key: "name", Value: req.Name})
	}

	var networkResourceList []types.NetworkResource

	networkResourceList, err = cli.NetworkList(context.TODO(), option)
	if err != nil {
		return nil, err
	}

	networkList := make([]model.Network, 0)
	for _, re := range networkResourceList {

		ipamDriver := make([]model.IPAMDriver, 0)

		for _, c := range re.IPAM.Config {
			ipamDriver = append(ipamDriver, model.IPAMDriver{Subnet: c.Subnet})
		}
		n := model.Network{
			Scope:   re.Scope,
			Id:      re.ID,
			Name:    re.Name,
			Driver:  re.Driver,
			Created: docker.TimestampToString(re.Created),
			IPAM: model.IPAM{
				Driver: re.IPAM.Driver,
			},
		}

		if len(ipamDriver) > 0 {
			n.IPAM.Config = ipamDriver[0]
		}

		networkList = append(networkList, n)
	}

	sort.Slice(networkList, func(i, j int) bool {
		if docker.StringToTimestamp(networkList[i].Created) > docker.StringToTimestamp(networkList[j].Created) {
			return true
		} else {
			return false
		}
	})

	return networkList, nil
}

// RemoveNetwork 删除网络
func (e *NetworkService) RemoveNetwork(host string, ids []string) error {

	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return err
	}

	for _, id := range ids {
		err := cli.NetworkRemove(context.Background(), id)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateNetwork 创建网络
func (e *NetworkService) CreateNetwork(host string, n model.Network) error {
	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return err
	}

	configs := make([]network.IPAMConfig, 0)
	if n.IPAM.Config.Subnet != "" {
		configs = append(configs, network.IPAMConfig{
			Subnet:  n.IPAM.Config.Subnet,
			Gateway: n.IPAM.Config.Gateway,
			IPRange: n.IPAM.Config.IPRange,
		})
	}
	_, err = cli.NetworkCreate(context.Background(), n.Name, types.NetworkCreate{
		Driver: n.Driver,
		IPAM: &network.IPAM{
			Driver: n.IPAM.Driver,
			Config: configs,
		}})
	if err != nil {
		return err
	}

	return nil
}
