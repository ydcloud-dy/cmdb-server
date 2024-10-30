package model

// SearchImage 搜索镜像
type SearchImage struct {
	Name     string `json:"name" form:"name"`
	OrderKey string `json:"orderKey" form:"orderKey"` // 排序
	Desc     bool   `json:"desc" form:"desc"`         // 排序方式:升序false(默认)|降序true
	Pagination
}

type SearchImageRes struct {
	Pagination
	Items []Image `json:"items"`
}

// Image 镜像
type Image struct {
	Tag     string `json:"tag"`
	Size    string `json:"size"`
	Id      string `json:"id"`
	Created string `json:"created"`
}

// PullImage 下载镜像参数
type PullImage struct {
	Name     string `json:"name"`
	PlatFrom string `json:"platFrom"`
	Auth     Auth   `json:"auth"`
}

// Auth 下载镜像认证参数
type Auth struct {
	Enable        bool   `json:"enable"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	ServerAddress string `json:"serverAddress"`
}

// RemoveImage 删除镜像
type RemoveImage struct {
	Ids []string `json:"Ids"`
}
