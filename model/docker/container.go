package model

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"

	"io"
	"strings"
)

// SearchContainer 搜索容器
type SearchContainer struct {
	Name     string `json:"name" form:"name"`
	OrderKey string `json:"orderKey" form:"orderKey"` // 排序
	Desc     bool   `json:"desc" form:"desc"`         // 排序方式:升序false(默认)|降序true
	State    string `json:"state" form:"state"`
	Pagination
}

type SearchContainerRes struct {
	Pagination
	Items []types.Container `json:"items"`
}

// Response 响应参数
type Response struct {
	code int    // 状态码
	msg  string // 返回信息
}

// AddContainer 创建容器
type AddContainer struct {
	ContainerConfig     ContainerConfig     `json:"containerConfig"`
	ContainerHostConfig ContainerHostConfig `json:"containerHostConfig"`
	ContainerNetwork    ContainerNetwork    `json:"containerNetwork"`
	Name                string              `json:"name"`
}

// ContainerConfig 容器配置
type ContainerConfig struct {
	Hostname     string  `json:"hostname"`
	DomainName   string  `json:"domainName"`
	Image        string  `json:"image"`
	User         string  `json:"user"`
	AttachStdin  bool    `json:"attachStdin"`
	AttachStdout bool    `json:"attachStdout"`
	AttachStderr bool    `json:"attachStderr"`
	Envs         []Env   `json:"envs"`
	Cmd          string  `json:"cmd"`
	Tty          bool    `json:"tty"`
	Entrypoint   string  `json:"entrypoint"`
	Labels       []Label `json:"labels"`
	OpenStdin    bool    `json:"openStdin"`
	WorkingDir   string  `json:"workingDir"`
	MacAddress   string  `json:"macAddress"`
	Console      string  `json:"console"`
}

// EnvString 环境变量转换字符串
func (c *ContainerConfig) EnvString() []string {
	r := make([]string, 0)
	for _, e := range c.Envs {
		r = append(r, e.Key+"="+e.Value)
	}
	return r
}

// CmdArray 命令转换数组
func (c *ContainerConfig) CmdArray() []string {
	if c.Cmd == "" {
		return []string{}
	}
	return strings.Split(c.Cmd, " ")
}

// EntrypointArray 命令转换数组
func (c *ContainerConfig) EntrypointArray() []string {
	if c.Entrypoint == "" {
		return []string{}
	}
	return strings.Split(c.Entrypoint, " ")
}

// LabelMap 标签转换
func (c *ContainerConfig) LabelMap() map[string]string {
	r := make(map[string]string, 0)
	for _, e := range c.Envs {
		r[e.Key] = e.Value
	}
	return r
}

type Env struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ContainerHostConfig 容器主机配置
type ContainerHostConfig struct {
	RestartPolicy  string           `json:"restartPolicy"`
	AutoRemove     bool             `json:"autoRemove"`
	LogConfig      LogConfig        `json:"logConfig"`
	ContainerMount []ContainerMount `json:"containerMount"`
	PortBinding    []PortBinding    `json:"portBinding"`
	ExtraHosts     []ExtraHost      `json:"extraHosts"`
	Privileged     bool             `json:"privileged"`
	ShmSize        int64            `json:"shmSize"`
	Sysctls        []Sysctl         `json:"sysctls"`
	Resource       Resource         `json:"resource"`
}

// Sysctl 系统配置
type Sysctl struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Resource 资源配置
type Resource struct {
	Memory            int64 `json:"memory"`
	MemoryReservation int64 `json:"memoryReservation"`
	NanoCpus          int64 `json:"nanoCpus"`
}

// ExtraHostsArr 容器主机配置
func (c *ContainerHostConfig) ExtraHostsArr() []string {
	r := make([]string, 0)
	for _, e := range c.ExtraHosts {
		r = append(r, e.Value)
	}
	return r
}

// SysctlsMap 容器主机配置
func (c *ContainerHostConfig) SysctlsMap() map[string]string {
	r := make(map[string]string, 0)
	for _, e := range c.Sysctls {
		r[e.Name] = e.Value
	}

	return r
}

// ExtraHost 容器主机扩展配置
type ExtraHost struct {
	Value string `json:"value,omitempty"`
}

// PortBinding 容器端口绑定
type PortBinding struct {
	Host      string `json:"host,omitempty"`
	Container string `json:"container"`
	Protocol  string `json:"protocol"`
}

// ContainerMount 容器挂载参数
type ContainerMount struct {
	Type        mount.Type `json:"type,omitempty"`
	Source      string     `json:"source"`
	Target      string     `json:"target"`
	StorageType string     `json:"storageType"`
}

// LogConfig 容器日志配置
type LogConfig struct {
	Type   string            `json:"type"`
	Config map[string]string `json:"config"`
}

// ContainerNetwork 容器网络配置
type ContainerNetwork struct {
	Name            string          `json:"name"`
	NetworkEndpoint NetworkEndpoint `json:"networkEndpoint"`
}

// NetworkEndpoint 网络配置
type NetworkEndpoint struct {
	Gateway    string            `json:"gateway"`
	IPAddress  string            `json:"ipAddress"`
	MacAddress string            `json:"macAddress"`
	DriverOpts map[string]string `json:"driverOpts"`
}

// EndpointIPAMConfig 网络endpoint配置
type EndpointIPAMConfig struct {
	IPv4Address  string   `json:"ipv4Address,omitempty"`
	IPv6Address  string   `json:"ipv6Address,omitempty"`
	LinkLocalIPs []string `json:"linkLocalIPs,omitempty"`
}

// UpdateContainer 更新容器配置
type UpdateContainer struct {
	Id     string                 `json:"Id"`
	Config container.UpdateConfig `json:"Config"`
}

// RemoveContainer 删除容器配置
type RemoveContainer struct {
	Ids           []string `json:"Ids"`
	RemoveVolumes bool     `json:"RemoveVolumes"`
	RemoveLinks   bool     `json:"RemoveLinks"`
	Force         bool     `json:"Force"`
}

// RestartContainer 容器重启配置
type RestartContainer struct {
	Id string `json:"Id"`
}

// StartContainer 启动容器配置
type StartContainer struct {
	Id string `json:"Id"`
}

// StopContainer 停止容器配置
type StopContainer struct {
	Id string `json:"Id"`
}

// StatsContainer 容器状态配置
type StatsContainer struct {
	Id     string `json:"Id"`
	Stream bool   `json:"Stream"`
}

// GetContainerLog 获取容器日志配置
type GetContainerLog struct {
	Datetimerange string `json:"datetimerange" form:"datetimerange"`
	ContainerId   string `json:"containerId" form:"containerId" `
	Since         string `json:"since" form:"since"`
	Stdout        bool   `json:"stdout" form:"stdout"`
	Stderr        bool   `json:"stderr" form:"stderr"`
	Tail          string `json:"tail" form:"tail"`
	Timestamps    bool   `json:"timestamps" form:"timestamps"`
	Text          string `json:"text" form:"text"`
	Regexp        bool   `json:"regexp" form:"regexp"`
	UpperLower    bool   `json:"upperLower" form:"upperLower"`
}

// ContainerLogResponse 容器日志响应
type ContainerLogResponse struct {
	Reader io.ReadCloser
}

// ExecContainer 执行容器配置
type ExecContainer struct {
	Id       string `json:"id"`
	ExecType string `json:"execType"`
}

// ExecContainerResize 执行容器配置
type ExecContainerResize struct {
	Id     string `json:"id,omitempty"`
	Width  uint   `json:"width,omitempty"`
	Height uint   `json:"height,omitempty"`
}

// StatsContainerStatsRes 容器状态响应
type StatsContainerStatsRes struct {
	Time     string  `json:"time,omitempty"`
	CpuUsage float64 `json:"cpuUsage"`
	MemUsage float64 `json:"memUsage"`
	MemCache float64 `json:"memCache"`
	IORead   float64 `json:"ioRead"`
	IOWrite  float64 `json:"ioWrite"`
	RxBytes  float64 `json:"rxBytes"`
	TxBytes  float64 `json:"txBytes"`
}

// InspectContainer 获取容器详细信息
type InspectContainer struct {
	Id string `json:"id,omitempty" form:"id"`
}
