package model

import (
	"fmt"
)

// SearchVolume 搜索卷参数
type SearchVolume struct {
	Name     string `json:"name" form:"name"`
	OrderKey string `json:"orderKey" form:"orderKey"` // 排序
	Desc     bool   `json:"desc" form:"desc"`         // 排序方式:升序false(默认)|降序true
	Pagination
}

type SearchVolumeRes struct {
	Pagination
	Items []Volume `json:"items"`
}

// RemoveVolume 删除卷参数
type RemoveVolume struct {
	Ids []string `json:"Ids"`
}

// Volume 卷
type Volume struct {
	Name       string                 `json:"name" yaml:"name"`
	Driver     string                 `yaml:"driver" json:"driver"`
	MountPoint string                 `yaml:"mountPoint" json:"mountPoint"`
	Created    string                 `yaml:"created" json:"created"`
	Scope      string                 `yaml:"scope" json:"scope"`
	Labels     map[string]string      `json:"labels"`
	Status     map[string]interface{} `json:"status,omitempty"`
	Cifs       CifsOption             `json:"cifs"`
	Nfs        NfsOption              `json:"nfs"`
	Options    map[string]string      `json:"options"`
}

// CifsOption cifs配置
type CifsOption struct {
	Enable   bool   `json:"enable"`
	Addr     string `json:"addr"`
	Version  string `json:"version"`
	Device   string `json:"device"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *CifsOption) Option() string {
	return fmt.Sprintf("addr=%v,username=%v,password=%v,vers=%v", c.Addr, c.Username, c.Password, c.Version)
}

// NfsOption nfs配置
type NfsOption struct {
	Addr    string `json:"addr"`
	Device  string `json:"device"`
	Version string `json:"version"`
	Args    string `json:"args"`
	Enable  bool   `json:"enable"`
}

func (n *NfsOption) Option() string {
	return fmt.Sprintf("addr=%v,%v,vers=%v", n.Addr, n.Args, n.Version)
}
