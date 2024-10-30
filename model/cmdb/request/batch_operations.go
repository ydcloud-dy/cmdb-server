package request

import "time"

// 执行命令请求的结构体
type ExecuteRequest struct {
	Hosts    []string `json:"hosts" form:"hosts"`       // 主机列表
	Users    []string `json:"users" form:"users"`       //
	Ports    []int    `json:"ports" form:"ports"`       //
	UserId   uint     `json:"userId" form:"userId"`     //
	Commands []string `json:"commands" form:"commands"` // 要执行的命令列表
	Language string   `json:"language" form:"language"` // "shell" 或者"python"
}

// 每台主机执行结果
type HostExecResult struct {
	Host   string `json:"host" form:"host"`
	Output string `json:"output" form:"output"`
	Error  string `json:"error,omitempty" form:"error"`
}

// 响应体的结构体
type ExecuteResponse struct {
	AllHosts      []string         `json:"allHosts" form:"allHosts"`           // 所有的主机列表
	SuccessHosts  []string         `json:"successHosts" form:"successHosts"`   // 成功的主机列表
	FailureHosts  []string         `json:"failureHosts" form:"failureHosts"`   // 失败的主机列表
	ExecutionLogs []HostExecResult `json:"executionLogs" form:"executionLogs"` // 执行日志
	Status        string           `json:"status"`
}

type CommandExecutionLog struct {
	ID            uint      `gorm:"primaryKey"`
	UserId        uint      `gorm:"type:bigint"`
	Command       string    `gorm:"type:text"` // 执行的命令
	AllHosts      string    `gorm:"type:text"` // 所有主机
	SuccessHosts  string    `gorm:"type:text"` // 成功的主机
	FailureHosts  string    `gorm:"type:text"` // 失败的主机
	ExecutionLogs string    `gorm:"type:text"` // 执行的日志（序列化后的 JSON 字符串）
	Status        string    `gorm:"size:50"`   // 状态：成功或失败
	CreatedAt     time.Time // 自动保存创建时间
}

func (c *CommandExecutionLog) TableName() string {
	return "cmdb_execution_logs"
}
