package docker

import "fmt"

// UnitFormat 单位格式化
func UnitFormat(u int64) string {
	if u < 1024*1024 { // kb
		return fmt.Sprintf("%.1fKB", float64(u)/1024)
	} else if u < 1024*1024*1024 { // mb
		return fmt.Sprintf("%.1fMB", float64(u)/1024/1024)
	} else if u < 1024*1024*1024*1024 { // gb
		return fmt.Sprintf("%.1fGB", float64(u)/1024/1024/1024)
	}
	return fmt.Sprintf("%v", u)
}
