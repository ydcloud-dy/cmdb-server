package metrics

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type MetricsCategory struct {
	ClusterId uint   `json:"cluster_id"`
	Category  string `json:"category"`
	Nodes     string `json:"nodes,omitempty"`
	PVC       string `json:"pvc,omitempty"`
	Pods      string `json:"pods,omitempty"`
	Ingress   string `json:"ingress,omitempty"`
	Selector  string `json:"selector,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Status    string `json:"status,omitempty"`
	Start     int64  `json:"start,omitempty"`
	End       int64  `json:"end,omitempty"`
}

func bytesSent(ingress string, namespace string, statuses string) string {
	return fmt.Sprintf(`sum(rate(nginx_ingress_controller_bytes_sent_sum{ingress="%s",namespace="%s",status=~"%s"}[1m])) by (ingress, namespace)`, ingress, namespace, statuses)
}

func (mc *MetricsCategory) GenerateQuery() *PrometheusQuery {
	switch mc.Category {
	case "cluster":
		return &PrometheusQuery{
			MemoryUsage:               strings.Replace("sum(node_memory_MemTotal_bytes - (node_memory_MemFree_bytes + node_memory_Buffers_bytes + node_memory_Cached_bytes)) by (instance)", "_bytes", fmt.Sprintf("_bytes{instance=~\"%s\"}", mc.Nodes), -1),
			MemoryRequests:            fmt.Sprintf(`sum(kube_pod_container_resource_requests{node=~"%s", resource="memory"}) by (component)`, mc.Nodes),
			MemoryLimits:              fmt.Sprintf(`sum(kube_pod_container_resource_limits{node=~"%s", resource="memory"}) by (component)`, mc.Nodes),
			MemoryCapacity:            fmt.Sprintf(`sum(kube_node_status_capacity{node=~"%s", resource="memory"}) by (component)`, mc.Nodes),
			MemoryAllocatableCapacity: fmt.Sprintf(`sum(kube_node_status_allocatable{node=~"%s", resource="memory"}) by (component)`, mc.Nodes),
			CpuUsage:                  fmt.Sprintf(`sum(rate(node_cpu_seconds_total{instance=~"%s", mode=~"user|system"}[1m]))`, mc.Nodes),
			CpuRequests:               fmt.Sprintf(`sum(kube_pod_container_resource_requests{node=~"%s", resource="cpu"}) by (component)`, mc.Nodes),
			CpuLimits:                 fmt.Sprintf(`sum(kube_pod_container_resource_limits{node=~"%s", resource="cpu"}) by (component)`, mc.Nodes),
			CpuCapacity:               fmt.Sprintf(`sum(kube_node_status_capacity{node=~"%s", resource="cpu"}) by (component)`, mc.Nodes),
			CpuAllocatableCapacity:    fmt.Sprintf(`sum(kube_node_status_allocatable{node=~"%s", resource="cpu"}) by (component)`, mc.Nodes),
			PodUsage:                  fmt.Sprintf(`sum({__name__=~"kubelet_running_pod_count|kubelet_running_pods", node=~"%s"})`, mc.Nodes),
			PodCapacity:               fmt.Sprintf(`sum(kube_node_status_capacity{node=~"%s", resource="pods"}) by (component)`, mc.Nodes),
			PodAllocatableCapacity:    fmt.Sprintf(`sum(kube_node_status_allocatable{node=~"%s", resource="pods"}) by (component)`, mc.Nodes),
			FsSize:                    fmt.Sprintf(`sum(node_filesystem_size_bytes{instance=~"%s", mountpoint="/"}) by (kubernetes_node)`, mc.Nodes),
			FsUsage:                   fmt.Sprintf(`sum(node_filesystem_size_bytes{instance=~"%s", mountpoint="/"} - node_filesystem_avail_bytes{instance=~"%s", mountpoint="/"}) by (kubernetes_node)`, mc.Nodes, mc.Nodes),
		}
	case "nodes":
		return &PrometheusQuery{
			MemoryUsage:            `sum(node_memory_MemTotal_bytes - (node_memory_MemFree_bytes + node_memory_Buffers_bytes + node_memory_Cached_bytes)) by (instance)`,
			MemoryCapacity:         `sum(kube_node_status_capacity{resource="memory"}) by (node)`,
			MemoryRequests:         fmt.Sprintf(`sum(kube_pod_container_resource_requests{node=~"%s", resource="memory"}) by (node)`, mc.Nodes),
			CpuUsage:               `sum(rate(node_cpu_seconds_total{mode=~"user|system"}[1m])) by (instance)`,
			CpuRequests:            fmt.Sprintf(`sum(kube_pod_container_resource_requests{node=~"%s", resource="cpu"}) by (node)`, mc.Nodes),
			CpuCapacity:            fmt.Sprintf(`sum(kube_node_status_capacity{resource="cpu", node=~"%s"}) by (node)`, mc.Nodes),
			FsSize:                 fmt.Sprintf(`sum(node_filesystem_size_bytes{mountpoint="/", instance=~"%s"}) by (instance)`, mc.Nodes),
			FsUsage:                `sum(node_filesystem_size_bytes{mountpoint="/"} - node_filesystem_avail_bytes{mountpoint="/"}) by (instance)`,
			PodUsage:               fmt.Sprintf(`sum({__name__=~"kubelet_running_pod_count|kubelet_running_pods", node=~"%s"})`, mc.Nodes),
			PodCapacity:            fmt.Sprintf(`sum(kube_node_status_capacity{node=~"%s", resource="pods"}) by (node)`, mc.Nodes),
			PodAllocatableCapacity: fmt.Sprintf(`sum(kube_node_status_allocatable{node=~"%s", resource="pods"}) by (component)`, mc.Nodes),
		}
	case "pods":
		return &PrometheusQuery{
			CpuUsage:              fmt.Sprintf(`sum(rate(container_cpu_usage_seconds_total{image!="",pod="%s", namespace="%s"}[1m])) by (pod,namespace)`, mc.Pods, mc.Namespace),
			CpuRequests:           fmt.Sprintf(`sum(kube_pod_container_resource_requests{resource="cpu", pod=~"%s", namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			CpuLimits:             fmt.Sprintf(`sum(kube_pod_container_resource_limits{resource="cpu", pod=~"%s", namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			MemoryUsage:           fmt.Sprintf(`sum(container_memory_working_set_bytes{image!="", pod=~"%s", namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			MemoryRequests:        fmt.Sprintf(`sum(kube_pod_container_resource_requests{resource="memory", pod=~"%s", namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			MemoryLimits:          fmt.Sprintf(`sum(kube_pod_container_resource_limits{resource="memory", pod=~"%s", namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			FsUsage:               fmt.Sprintf(`sum(container_fs_usage_bytes{container!="POD",container!="",pod=~"%s", namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			FsWrite:               fmt.Sprintf(`sum(container_fs_writes_bytes_total{container!="", pod=~"%s", namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			FsRead:                fmt.Sprintf(`sum(container_fs_reads_bytes_total{container!="", pod="%s", namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			NetworkReceive:        fmt.Sprintf(`sum(container_network_receive_bytes_total{pod=~"%s",namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			NetworkTransmit:       fmt.Sprintf(`sum(container_network_transmit_bytes_total{pod=~"%s",namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			PodTcpEstablishedConn: fmt.Sprintf(`sum(inspector_pod_tcpsummarytcpestablishedconn{target_pod=~"%s", target_namespace="%s"} ) by  (target_pod, target_namespace)`, mc.Pods, mc.Namespace),
			PodTcpTimewaitConn:    fmt.Sprintf(`sum(inspector_pod_tcpsummarytcptimewaitconn{target_pod=~"%s", target_namespace="%s"} ) by  (target_pod, target_namespace)`, mc.Pods, mc.Namespace),
		}
	case "ingress":
		return &PrometheusQuery{
			BytesSentSuccess:        bytesSent(mc.Ingress, mc.Namespace, "^2\\\\d*"),
			BytesSent3XX:            bytesSent(mc.Ingress, mc.Namespace, "^3\\\\d*"),
			BytesSent4XX:            bytesSent(mc.Ingress, mc.Namespace, "^4\\\\d*"),
			BytesSentFailure:        bytesSent(mc.Ingress, mc.Namespace, "^5\\\\d*"),
			RequestDurationSeconds:  fmt.Sprintf(`sum(rate(nginx_ingress_controller_request_duration_seconds_sum{ingress="%s",namespace="%s"}[1m])) by (ingress, namespace)`, mc.Ingress, mc.Namespace),
			ResponseDurationSeconds: fmt.Sprintf(`sum(rate(nginx_ingress_controller_response_duration_seconds_sum{ingress="%s",namespace="%s"}[1m])) by (ingress, namespace)`, mc.Ingress, mc.Namespace),
		}
	}
	return nil
}

type MetricsQuery struct {
	MemoryUsage               *MetricsCategory `json:"memoryUsage,omitempty"`
	MemoryRequests            *MetricsCategory `json:"memoryRequests,omitempty"`
	MemoryLimits              *MetricsCategory `json:"memoryLimits,omitempty"`
	MemoryCapacity            *MetricsCategory `json:"memoryCapacity,omitempty"`
	MemoryAllocatableCapacity *MetricsCategory `json:"memoryAllocatableCapacity,omitempty"`
	CpuUsage                  *MetricsCategory `json:"cpuUsage,omitempty"`
	CpuLimits                 *MetricsCategory `json:"cpuLimits,omitempty"`
	CpuRequests               *MetricsCategory `json:"cpuRequests,omitempty"`
	CpuCapacity               *MetricsCategory `json:"cpuCapacity,omitempty"`
	CpuAllocatableCapacity    *MetricsCategory `json:"cpuAllocatableCapacity,omitempty"`
	FsSize                    *MetricsCategory `json:"fsSize,omitempty"`
	FsUsage                   *MetricsCategory `json:"fsUsage,omitempty"`
	FsWrite                   *MetricsCategory `json:"fsWrite,omitempty"`
	FsRead                    *MetricsCategory `json:"fsRead,omitempty"`
	PodUsage                  *MetricsCategory `json:"podUsage,omitempty"`
	PodCapacity               *MetricsCategory `json:"podCapacity,omitempty"`
	PodAllocatableCapacity    *MetricsCategory `json:"podAllocatableCapacity,omitempty"`
	NetworkReceive            *MetricsCategory `json:"networkReceive,omitempty"`
	NetworkTransmit           *MetricsCategory `json:"networkTransmit,omitempty"`
	BytesSentSuccess          *MetricsCategory `json:"bytesSentSuccess,omitempty"`
	BytesSent3XX              *MetricsCategory `json:"bytesSent3XX,omitempty"`
	BytesSent4XX              *MetricsCategory `json:"bytesSent4XX,omitempty"`
	BytesSentFailure          *MetricsCategory `json:"bytesSentFailure,omitempty"`
	RequestDurationSeconds    *MetricsCategory `json:"requestDurationSeconds,omitempty"`
	ResponseDurationSeconds   *MetricsCategory `json:"responseDurationSeconds,omitempty"`
	WorkloadMemoryUsage       *MetricsCategory `json:"workloadMemoryUsage,omitempty"`
	PodTcpEstablishedConn     *MetricsCategory `json:"podTcpEstablishedConn,omitempty"`
	PodTcpTimewaitConn        *MetricsCategory `json:"podTcpTimewaitConn,omitempty"`
}

type PrometheusQuery struct {
	CpuUsage                  string
	CpuRequests               string
	CpuLimits                 string
	CpuCapacity               string
	WorkloadMemoryUsage       string
	CpuAllocatableCapacity    string
	MemoryUsage               string
	MemoryCapacity            string
	MemoryRequests            string
	MemoryLimits              string
	MemoryAllocatableCapacity string
	FsUsage                   string
	FsSize                    string
	FsWrite                   string
	FsRead                    string
	NetworkReceive            string
	NetworkTransmit           string
	PodUsage                  string
	PodCapacity               string
	PodAllocatableCapacity    string
	DiskUsage                 string
	DiskCapacity              string
	BytesSentSuccess          string
	BytesSent3XX              string
	BytesSent4XX              string
	BytesSentFailure          string
	RequestDurationSeconds    string
	ResponseDurationSeconds   string
	PodTcpEstablishedConn     string
	PodTcpTimewaitConn        string
}

func (pq *PrometheusQuery) GetValueByField(field string) string {
	e := reflect.ValueOf(pq).Elem()
	for i := 0; i < e.NumField(); i++ {
		if e.Type().Field(i).Name == field {
			return e.Field(i).Interface().(string)
		}
	}
	return ""
}

type PrometheusQueryResp struct {
	Status string                   `json:"status"`
	Data   *PrometheusQueryRespData `json:"data"`
}

type PrometheusQueryRespData struct {
	ResultType string                      `json:"resultType"`
	Result     []PrometheusQueryRespResult `json:"result"`
}

type PrometheusQueryRespResult struct {
	Metric interface{}   `json:"metric"`
	Values []interface{} `json:"values"`
}

type PrometheusTracker struct {
	// 添加读写锁来保护下面的 map
	sync.RWMutex
	Metrics map[string]*PrometheusQueryResp
}

func NewPrometheusTracker() *PrometheusTracker {
	return &PrometheusTracker{Metrics: map[string]*PrometheusQueryResp{}}
}

func (pt *PrometheusTracker) Get(key string) (*PrometheusQueryResp, bool) {
	pt.RLock()
	defer pt.RUnlock()
	val, ext := pt.Metrics[key]
	return val, ext
}

func (pt *PrometheusTracker) Set(key string, val *PrometheusQueryResp) {
	pt.Lock()
	defer pt.Unlock()
	pt.Metrics[key] = val
}
