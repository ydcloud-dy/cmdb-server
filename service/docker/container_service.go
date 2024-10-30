package service

import (
	global2 "DYCLOUD/global/docker"
	model "DYCLOUD/model/docker"
	docker "DYCLOUD/utils/docker/docker"
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/golang-module/carbon"
	"net"
	"sort"
	"strings"
)

type ContainerService struct {
	dockerClient *client.Client
}

// ListContainer 获取容器列表
func (e *ContainerService) ListContainer(host string, req model.SearchContainer) (res []types.Container, err error) {

	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return nil, err
	}

	res = make([]types.Container, 0)

	option := types.ContainerListOptions{
		All: true,
	}

	if strings.Trim(req.Name, " ") != "" {
		option.Filters = filters.NewArgs(filters.KeyValuePair{Key: "name", Value: req.Name})
	}

	res, err = cli.ContainerList(context.TODO(), option)

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].State == "running"
	})

	if err != nil {
		return make([]types.Container, 0), err
	}

	sort.Slice(res, func(i, j int) bool {
		if res[i].Created > res[j].Created {
			return true
		} else {
			return false
		}
	})

	tmpList := make([]types.Container, 0)
	if req.State != "all" {
		for _, v := range res {
			if req.State == v.State {
				tmpList = append(tmpList, v)
			}
		}
		return tmpList, nil
	}
	return res, nil
}

// AddContainer 创建容器
func (e *ContainerService) AddContainer(host string, c model.AddContainer) error {

	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return err
	}

	// 容器基础配置
	containerConf := &container.Config{
		Hostname:     c.ContainerConfig.Hostname,
		User:         c.ContainerConfig.User,
		Domainname:   c.ContainerConfig.DomainName,
		Image:        c.ContainerConfig.Image,
		AttachStdin:  c.ContainerConfig.AttachStdin,
		AttachStdout: c.ContainerConfig.AttachStdout,
		AttachStderr: c.ContainerConfig.AttachStderr,
		Env:          c.ContainerConfig.EnvString(),
		Labels:       c.ContainerConfig.LabelMap(),
		WorkingDir:   c.ContainerConfig.WorkingDir,
		MacAddress:   c.ContainerConfig.MacAddress,
	}
	cmdArr := c.ContainerConfig.CmdArray()
	if len(cmdArr) > 0 {
		containerConf.Cmd = cmdArr
	}

	entrypointArr := c.ContainerConfig.EntrypointArray()
	if len(entrypointArr) > 0 {
		containerConf.Entrypoint = entrypointArr
	}

	// 控制台配置
	if c.ContainerConfig.Console == "1" {
		containerConf.OpenStdin = true
		containerConf.Tty = true
	} else if c.ContainerConfig.Console == "2" {
		containerConf.Tty = true
	} else if c.ContainerConfig.Console == "3" {
		containerConf.OpenStdin = true
	}

	binds := make([]string, 0)
	mounts := make([]mount.Mount, 0)
	for _, m := range c.ContainerHostConfig.ContainerMount {

		mt := mount.Mount{
			Type:   m.Type,
			Source: m.Source,
			Target: m.Target,
		}
		if m.StorageType == "2" {
			mt.ReadOnly = true
		}
		mounts = append(mounts, mt)
		binds = append(binds, fmt.Sprintf("%v:%v", m.Source, m.Target))
	}

	var restartPolicy string
	switch c.ContainerHostConfig.RestartPolicy {
	case "1":
		restartPolicy = "no"
	case "2":
		restartPolicy = "always"
	case "3":
		restartPolicy = "on-failure"
	case "4":
		restartPolicy = "unless-stopped"
	}

	// 容器主机配置
	hostConf := &container.HostConfig{
		//Binds:      binds,
		AutoRemove: c.ContainerHostConfig.AutoRemove,

		RestartPolicy: container.RestartPolicy{
			Name: restartPolicy,
		},
		NetworkMode:  container.NetworkMode(c.ContainerNetwork.Name),
		Mounts:       mounts,
		ExtraHosts:   c.ContainerHostConfig.ExtraHostsArr(),
		PortBindings: make(nat.PortMap, 0),
		Privileged:   c.ContainerHostConfig.Privileged,
		ShmSize:      c.ContainerHostConfig.ShmSize * 1024 * 1024,
		Sysctls:      c.ContainerHostConfig.SysctlsMap(),
		Resources: container.Resources{
			Memory:            c.ContainerHostConfig.Resource.Memory * 1024 * 1024,
			NanoCPUs:          c.ContainerHostConfig.Resource.NanoCpus * 100000000,
			MemoryReservation: c.ContainerHostConfig.Resource.MemoryReservation * 1024 * 1024,
		},
	}

	if c.ContainerHostConfig.LogConfig.Type != "default" {
		hostConf.LogConfig = container.LogConfig{
			Type:   c.ContainerHostConfig.LogConfig.Type,
			Config: c.ContainerHostConfig.LogConfig.Config,
		}
	}

	for _, p := range c.ContainerHostConfig.PortBinding {
		hostPortList := docker.PortSection(p.Host)
		containerPortList := docker.PortSection(p.Container)

		if len(hostPortList) != len(containerPortList) {
			return errors.New("主机端口和容器端口数量不匹配")
		}

		for i := 0; i < len(hostPortList); i++ {
			hostConf.PortBindings[nat.Port(containerPortList[i]+"/"+p.Protocol)] = []nat.PortBinding{
				{
					HostPort: hostPortList[i],
				},
			}
		}
	}

	// 容器网络配置
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			c.ContainerNetwork.Name: {
				IPAMConfig: &network.EndpointIPAMConfig{
					IPv4Address: c.ContainerNetwork.NetworkEndpoint.IPAddress,
				},
				MacAddress: c.ContainerNetwork.NetworkEndpoint.MacAddress,
				DriverOpts: c.ContainerNetwork.NetworkEndpoint.DriverOpts,
			},
		},
	}

	response, err := cli.ContainerCreate(context.TODO(), containerConf, hostConf, networkConfig, nil, c.Name)
	if err != nil {
		return err
	}

	err = cli.ContainerStart(context.TODO(), response.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return nil
}

// UpdateContainer 修改容器配置
func (e *ContainerService) UpdateContainer(host string, c model.UpdateContainer) error {

	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return err
	}

	_, err = cli.ContainerUpdate(context.TODO(), c.Id, c.Config)
	if err != nil {
		return err
	}
	return nil
}

// RemoveContainer 删除容器
func (e *ContainerService) RemoveContainer(host string, c model.RemoveContainer) error {

	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return err
	}

	for _, id := range c.Ids {
		err := cli.ContainerRemove(context.TODO(), id, types.ContainerRemoveOptions{
			RemoveVolumes: c.RemoveVolumes,
			RemoveLinks:   c.RemoveLinks,
			Force:         c.Force,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// GetContainerLog 获取容器日志
func (e *ContainerService) GetContainerLog(host string, c *model.GetContainerLog) (model.ContainerLogResponse, error) {

	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return model.ContainerLogResponse{}, err
	}

	logOptions := types.ContainerLogsOptions{
		ShowStdout: c.Stdout,
		ShowStderr: c.Stderr,
		Tail:       c.Tail,
		Timestamps: c.Timestamps,
	}

	if c.Datetimerange != "" {
		rangeList := strings.Split(c.Datetimerange, ",")
		if len(rangeList) > 1 {
			startDate := rangeList[0]
			endDate := rangeList[1]
			logOptions.Since = fmt.Sprintf("%v", carbon.SetTimezone(carbon.UTC).Parse(startDate).Timestamp())
			logOptions.Until = fmt.Sprintf("%v", carbon.SetTimezone(carbon.UTC).Parse(endDate).Timestamp())
		}
	}

	logs, err := cli.ContainerLogs(context.TODO(), c.ContainerId, logOptions)
	return model.ContainerLogResponse{
		Reader: logs,
	}, err
}

// RestartContainer 重启容器
func (e *ContainerService) RestartContainer(host string, c model.RestartContainer) error {

	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return err
	}

	timeout := 60
	err = cli.ContainerRestart(context.TODO(), c.Id, container.StopOptions{Timeout: &timeout})
	if err != nil {
		return err
	}

	return nil
}

// StartContainer 启动容器
func (e *ContainerService) StartContainer(host string, c model.StartContainer) error {
	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return err
	}

	err = cli.ContainerStart(context.TODO(), c.Id, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return nil
}

// StopContainer 停止容器
func (e *ContainerService) StopContainer(host string, c model.StopContainer) error {

	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return err
	}

	timeout := 30
	err = cli.ContainerStop(context.TODO(), c.Id, container.StopOptions{Timeout: &timeout})
	if err != nil {
		return err
	}

	return nil
}

// GetContainerStats  获取容器资源统计
func (e *ContainerService) GetContainerStats(host string, c model.StatsContainer) (model.StatsContainerStatsRes, error) {

	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return model.StatsContainerStatsRes{}, err
	}

	stats, err := cli.ContainerStats(context.TODO(), c.Id, false)
	if err != nil {
		return model.StatsContainerStatsRes{}, err
	}

	defer stats.Body.Close()

	var statsData types.StatsJSON
	if err := json.NewDecoder(stats.Body).Decode(&statsData); err != nil {
		return model.StatsContainerStatsRes{}, err
	}

	cpuPer := docker.CalculateCPUPercent(statsData)
	memoryUsage, memoryCache := docker.CalculateMemoryUsage(statsData)
	ioRead, ioWrite := docker.CalculateIOUsage(statsData)
	rx, tx := docker.CalculateNetworkUsage(statsData)
	res := model.StatsContainerStatsRes{
		Time:     statsData.Read.Format("15:04:05"),
		CpuUsage: cpuPer,
		MemUsage: memoryUsage,
		MemCache: memoryCache,
		IORead:   ioRead,
		IOWrite:  ioWrite,
		RxBytes:  rx,
		TxBytes:  tx,
	}

	return res, nil
}

// ExecContainer 重启容器
func (e *ContainerService) ExecContainer(host string, c model.ExecContainer) (types.IDResponse, error) {

	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return types.IDResponse{}, err
	}

	idResponse, err := cli.ContainerExecCreate(context.Background(), c.Id, types.ExecConfig{
		AttachStderr: true,
		AttachStdin:  true,
		AttachStdout: true,
		Cmd:          []string{c.ExecType},
		Tty:          true,
	})
	if err != nil {
		return types.IDResponse{}, err
	}

	return idResponse, nil
}

// ExecAttachContainer 容器终端服务
func (e *ContainerService) ExecAttachContainer(host string, id string) (net.Conn, *bufio.Reader, error) {
	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return nil, nil, err
	}
	attach, err := cli.ContainerExecAttach(context.Background(), id, types.ExecStartCheck{
		Tty: true,
	})
	return attach.Conn, attach.Reader, err
}

// ExecContainerResize 修改container 终端大小
func (e *ContainerService) ExecContainerResize(host string, c model.ExecContainerResize) error {
	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return err
	}
	err = cli.ContainerExecResize(context.Background(), c.Id, types.ResizeOptions{
		Width:  c.Width,
		Height: c.Height,
	})
	return err
}

// InspectContainer 获取容器详细信息
func (e *ContainerService) InspectContainer(host string, c model.InspectContainer) (types.ContainerJSON, error) {
	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return types.ContainerJSON{}, err
	}
	return cli.ContainerInspect(context.Background(), c.Id)
}
