package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
)

// PortSection 端口段格式化
func PortSection(port string) []string {
	portList := make([]string, 0)
	if strings.Index(port, "-") == -1 {
		portList = append(portList, port)
		return portList
	}

	r := strings.Split(port, "-")
	var err error
	var begin, end int
	begin, err = strconv.Atoi(r[0])
	if err != nil {
		return portList
	}
	end, err = strconv.Atoi(r[1])
	if err != nil {
		return portList
	}

	for i := begin; i <= end; i++ {
		portList = append(portList, strconv.Itoa(i))
	}

	return portList
}

// CalculateCPUPercent 计算CPU使用率
func CalculateCPUPercent(stats types.StatsJSON) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)

	cpuPercent := 0.0
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}

	res, _ := decimal.NewFromFloat(cpuPercent).RoundBank(1).Float64()
	return res
}

// CalculateMemoryUsage 计算内存使用量
func CalculateMemoryUsage(stats types.StatsJSON) (float64, float64) {
	memoryUsageVal, _ := decimal.NewFromFloat(float64(stats.Stats.MemoryStats.Usage)).RoundBank(1).Float64()
	memoryCacheVal, _ := decimal.NewFromFloat(float64(stats.Stats.MemoryStats.Stats["cache"])).RoundBank(1).Float64()
	return memoryUsageVal, memoryCacheVal
}

// CalculateNetworkUsage 计算网络使用量
func CalculateNetworkUsage(stats types.StatsJSON) (float64, float64) {

	var inputTraffic, outputTraffic float64
	for _, network := range stats.Networks {
		inputTraffic += float64(network.RxBytes)
		outputTraffic += float64(network.TxBytes)
	}

	inputTrafficVal, _ := decimal.NewFromFloat(inputTraffic).RoundBank(1).Float64()
	outputTrafficVal, _ := decimal.NewFromFloat(outputTraffic).RoundBank(1).Float64()
	return inputTrafficVal, outputTrafficVal
}

// CalculateIOUsage 计算磁盘io使用量
func CalculateIOUsage(stats types.StatsJSON) (float64, float64) {
	var diskRead, diskWrite float64
	for _, disk := range stats.BlkioStats.IoServiceBytesRecursive {
		if disk.Op == "Read" {
			diskRead += float64(disk.Value)
		} else if disk.Op == "Write" {
			diskWrite += float64(disk.Value)
		}
	}
	diskReadVal, _ := decimal.NewFromFloat(diskRead).RoundBank(1).Float64()
	diskWriteVal, _ := decimal.NewFromFloat(diskWrite).RoundBank(1).Float64()
	return diskReadVal, diskWriteVal
}
