package huawei

var (
	ECSStatus = map[string]string{
		"BUILD":             "创建中",
		"REBOOT":            "重启",
		"HARD_REBOOT":       "强制重启",
		"REBUILD":           "重建中",
		"MIGRATING":         "热迁移中",
		"RESIZE":            "开始变更",
		"ACTIVE":            "运行中",
		"SHUTOFF":           "正常停止",
		"REVERT_RESIZE":     "回退变更规格",
		"VERIFY_RESIZE":     "校验变更完成配置",
		"ERROR":             "异常状态",
		"DELETED":           "正常删除",
		"SHELVED":           "搁置状态",
		"SHELVED_OFFLOADED": "卷启动的实例处于搁置状态",
		"UNKNOWN":           "未知状态",
	}

	LoadBalancerStatus = map[string]string{
		"ACTIVE":         "运行中",
		"PENDING_DELETE": "删除中",
	}

	RdsStatus = map[string]string{
		"BUILD":                   "创建中",
		"ACTIVE":                  "运行中",
		"FAILED":                  "异常",
		"FROZEN":                  "冻结",
		"MODIFYING":               "扩容中",
		"REBOOTING":               "重启中",
		"RESTORING":               "恢复中",
		"MODIFYING INSTANCE TYPE": "转主备中",
		"SWITCHOVER":              "主备切换中",
		"MIGRATING":               "迁移中",
		"BACKING UP":              "备份中",
		"MODIFYING DATABASE PORT": "修改数据库端口中",
		"STORAGE FULL":            "磁盘空间满",
	}
)
