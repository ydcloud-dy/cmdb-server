package aliyun

var (
	RdsStatus = map[string]string{
		"Creating":                  "创建中",
		"Running":                   "运行中",
		"Deleting":                  "删除中",
		"Rebooting":                 "重启中",
		"DBInstanceClassChanging":   "升降级中",
		"TRANSING":                  "迁移中",
		"EngineVersionUpgrading":    "迁移版本中",
		"TransingToOthers":          "迁移数据到其他RDS中",
		"GuardDBInstanceCreating":   "生产灾备实例中",
		"Restoring":                 "备份恢复中",
		"Importing":                 "数据导入中",
		"ImportingFromOthers":       "从其他RDS实例导入数据中",
		"DBInstanceNetTypeChanging": "内外网切换中",
		"GuardSwitching":            "容灾切换中",
		"INS_CLONING":               "实例克隆中",
		"Released":                  "已释放实例",
	}

	ECSStatus = map[string]string{
		"Pending":  "创建中",
		"Running":  "运行中",
		"Starting": "启动中",
		"Stopping": "停止中",
		"Stopped":  "已停止",
	}

	LoadBalancerStatus = map[string]string{
		"inactive": "已停止",
		"active":   "运行中",
		"locked":   "已锁定",
	}
)
