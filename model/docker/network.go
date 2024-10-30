package model

// SearchNetwork 网络搜索参数
type SearchNetwork struct {
	Name     string `json:"name" form:"name"`
	OrderKey string `json:"orderKey" form:"orderKey"` // 排序
	Desc     bool   `json:"desc" form:"desc"`         // 排序方式:升序false(默认)|降序true
	Pagination
}

// RemoveNetwork 网络删除参数
type RemoveNetwork struct {
	Ids []string `json:"Ids"`
}

type SearchNetworkRes struct {
	Pagination
	Items []Network `json:"items"`
}

// Network 网络
type Network struct {
	Name    string `yaml:"name" json:"name"`
	Id      string `yaml:"id" json:"id"`
	Scope   string `yaml:"scope" json:"scope"`
	Created string `yaml:"created" json:"created"`
	Driver  string `yaml:"driver" json:"driver"`
	IPAM    IPAM   `yaml:"ipam" json:"ipam"`
}

// IPAM 网络ipam配置
type IPAM struct {
	Driver string     `yaml:"driver" json:"driver"`
	Config IPAMDriver `yaml:"config" json:"config" `
}

// IPAMDriver 网络ipam驱动配置
type IPAMDriver struct {
	Subnet  string `yaml:"subnet" json:"subnet"`
	Gateway string `json:"gateway" yaml:"gateway"`
	IPRange string `json:"ipRange" yaml:"ipRange"`
}
