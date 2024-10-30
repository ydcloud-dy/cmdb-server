package cmdb

import "DYCLOUD/global"

// 主机管理 结构体  AssetHosts
type CmdbHosts struct {
	global.DYCLOUD_MODEL
	Name       string `json:"name" form:"name" gorm:"column:name;comment:;size:255;"` //name字段
	Port       *int   `json:"port" form:"port" gorm:"column:port;comment:;size:10;"`  //port字段
	Project    int    `json:"project" form:"project" gorm:"column:project"`
	Username   string `json:"username" form:"username" gorm:"column:username;comment:;size:255;"` //username字段
	ServerHost string `json:"serverHost" form:"serverHost" gorm:"column:serverHost"`
	PrivateKey string `json:"privateKey" form:"privateKey" gorm:"column:privateKey;type:text"`
	Password   string `json:"password" form:"password" gorm:"column:password"`
	Labels     string `json:"labels" form:"labels" gorm:"column:labels;comment:;size:255;"`       //labels字段
	Note       string `json:"note" form:"note" gorm:"column:note;comment:;size:255;"`             //note字段
	CpuModel   string `json:"cpuModel" form:"cpuModel" gorm:"column:cpuModel;comment:;size:255;"` //cpuModel字段
	CpuCount   string `json:"cpuCount" form:"cpuCount" gorm:"column:cpuCount;comment:;size:255;"` //cpuCount字段
	//CpuVcpus     string `json:"cpuVcpus" form:"cpuVcpus" gorm:"column:cpuVcpus;comment:;size:255;"`             //cpuVcpus字段
	Memory    string `json:"memory" form:"memory" gorm:"column:memory;comment:;size:255;"`          //memory字段
	DiskTotal string `json:"diskTotal" form:"diskTotal" gorm:"column:diskTotal;comment:;size:255;"` //diskTotal字段
	DiskInfo  string `json:"diskInfo" form:"diskInfo" gorm:"column:diskInfo;comment:;size:255;"`    //diskInfo字段
	Os        string `json:"os" form:"os" gorm:"column:os;comment:;size:255;"`                      //os字段
	OsVersion string `json:"osVersion" form:"osVersion" gorm:"column:osVersion;comment:;size:255;"` //osVersion字段
	OsArch    string `json:"osArch" form:"osArch" gorm:"column:osArch;comment:;size:255;"`          //osArch字段
	//HardwareInfo string `json:"hardwareInfo" form:"hardwareInfo" gorm:"column:hardwareInfo;comment:;size:255;"` //hardwareInfo字段
	Status    string `json:"status" form:"status" gorm:"column:status;comment:;size:255;"`          //status字段
	PublicIP  string `json:"publicIP" form:"publicIP" gorm:"column:publicIP;comment:;size:255;"`    // 公网ip
	PrivateIP string `json:"privateIP" form:"privateIP" gorm:"column:privateIP;comment:;size:255;"` // 公网ip

	CreatedBy uint `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint `gorm:"column:deleted_by;comment:删除者"`
}

// TableName cmdbHosts表 CmdbHosts自定义表名 cmdb_Hosts
func (CmdbHosts) TableName() string {
	return "cmdb_Hosts"
}
