package docker

import (
	"time"
)

// TimeFormat 时间戳格式化日期字符串
func TimeFormat(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
}

// TimestampToString 时间戳转换日期
func TimestampToString(d time.Time) string {
	return d.Format("2006-01-02 15:04:05")
}

// StringTimeFormat 日期转换时间戳
func StringTimeFormat(d string) string {
	p, err := time.Parse(time.RFC3339, d)
	if err != nil {
	}
	return p.Format("2006-01-02 15:04:05")
}

// StringToTimestamp 日期转换时间戳
func StringToTimestamp(val string) int64 {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", val, time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}
