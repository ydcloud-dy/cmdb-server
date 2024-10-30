package tencent

var (
	ECSStatus = map[string]string{
		"RUNNING":       "运行中",
		"PENDING":       "创建中",
		"LAUNCH_FAILED": "创建失败",
		"STARTING":      "启动中",
		"STOPPING":      "关机中",
		"REBOOTING":     "重启中",
		"SHUTDOWN":      "停止待销毁",
		"TERMINATING":   "销毁中",
	}

	LoadBalancerStatus = map[uint64]string{
		0: "运行中",
		1: "创建中",
	}

	RdsStatus = map[int64]string{
		0: "创建中",
		1: "运行中",
		4: "正在进行隔离操作",
		5: "已隔离",
	}
)
