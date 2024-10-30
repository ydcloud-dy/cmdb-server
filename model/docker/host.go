// 自动生成模板Host
package model

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/request"
	"time"
)

const (
	HostTypeSocket = iota
	HostTypeApi
)

// 主机列表 结构体  Host
type Host struct {
	global.DYCLOUD_MODEL
	SocketPath  string `json:"socketPath" form:"socketPath" gorm:"column:socketPath;comment:Socket文件路径;" ` //主机名称
	Type        int    `json:"type" form:"type" gorm:"column:type;comment:主机连接类型;" `
	Name        string `json:"name" form:"name" gorm:"column:name;comment:主机名称;" binding:"required"`                              //主机名称
	Description string `json:"description" form:"description" gorm:"column:description;comment:主机描述;"`                            //主机描述
	EnableTls   *bool  `json:"enableTls" form:"enableTls" gorm:"default:0;column:enable_tls;comment:是否开启tls;" binding:"required"` //是否开启tls
	TlsCa       string `json:"tlsCa" form:"tlsCa" gorm:"column:tls_ca;type:text;comment:tls ca证书;"`                               //tls ca证书
	TlsCert     string `json:"tlsCert" form:"tlsCert" gorm:"column:tls_cert;type:text;comment:tls crt证书;"`                        //tls crt证书
	TlsKey      string `json:"tlsKey" form:"tlsKey" gorm:"column:tls_key;type:text;comment:tls key证书;"`                           //tls key证书
	SkipCert    *bool  `json:"skipCert" form:"skipCert" gorm:"column:skip_cert;comment:是否跳过安全证书检查;"`                              //是否跳过安全证书检查
	ApiAddress  string `json:"apiAddress" form:"apiAddress" gorm:"column:api_address;comment:DockerApi访问地址;"`                     //DockerApi访问地址
	Port        *int   `json:"port" form:"port" gorm:"column:port;comment:docker api访问端口;"`                                       //docker api访问端口
	CreatedBy   uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy   uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy   uint   `gorm:"column:deleted_by;comment:删除者"`
}

// TableName 主机列表 Host自定义表名 host
func (Host) TableName() string {
	return "host"
}

type HostSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`

	request.PageInfo
}

type HostGenTlsScript struct {
	Script string `json:"script"`
}

type HostCheckResponse struct {
	Success bool `json:"success"`
}
