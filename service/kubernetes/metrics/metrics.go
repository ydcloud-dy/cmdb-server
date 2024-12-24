package metrics

import (
	"DYCLOUD/global"
	cluster2 "DYCLOUD/model/kubernetes/cluster"
	"DYCLOUD/model/kubernetes/metrics"
	"DYCLOUD/service/kubernetes/cluster"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"
)

type MetricsService struct{}

var (
	clusterService = cluster.K8sClusterService{}
)

// 认证
func (m MetricsService) Auth(clusterObj cluster2.K8sCluster) (httpClient http.Client) {
	httpClient = http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	if clusterObj.PrometheusAuthType == 2 {
		httpClient = http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				Proxy: func(req *http.Request) (*url.URL, error) {
					req.SetBasicAuth(clusterObj.PrometheusUser, clusterObj.PrometheusPwd)
					return nil, nil
				},
			}}
	}
	return httpClient
}

// 健康检查
func (m MetricsService) Health(clusterId int) (url string, httpClient http.Client, err error) {
	clusterObj, err := clusterService.GetK8sCluster(clusterId)

	if err != nil {
		return url, httpClient, err
	}

	// 超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()

	// 生成httpClient
	httpClient = m.Auth(clusterObj)

	prourl := fmt.Sprintf("%s/-/ready", clusterObj.PrometheusUrl)
	readyReq, err := http.NewRequest("GET", prourl, nil)
	if err != nil {
		global.DYCLOUD_LOG.Error("prometheus service ready request failed", zap.Any("err", err))
		return
	}

	readyResp, err := httpClient.Do(readyReq.WithContext(ctx))
	if err != nil {
		global.DYCLOUD_LOG.Error("prometheus service ready request failed", zap.Any("err", err))
		return
	}

	// 如果还没有ready，则直接返回前端空数据
	if readyResp.StatusCode != http.StatusOK {
		global.DYCLOUD_LOG.Error("prometheus resp no ready", zap.Any("err", err))
		return
	}

	return clusterObj.PrometheusUrl, httpClient, err
}

//@function: GetMetrics
//@description: 普罗米修斯监控数据获取
//@param: requestInfo kubernetesReq.ResourceParamRequest, queryinfo kubernetesReq.ResourceCreateRequest
//@return: err error

func (m MetricsService) GetMetrics(mt metrics.MetricsQuery) (t map[string]*metrics.PrometheusQueryResp, err error) {
	step := 60
	end := time.Now().Unix()
	start := end - 3600

	// tracker
	tracker := metrics.NewPrometheusTracker()
	wg := sync.WaitGroup{}
	e := reflect.ValueOf(&mt).Elem()
	for i := 0; i < e.NumField(); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fName := e.Type().Field(i).Name
			fValue := e.Field(i).Interface().(*metrics.MetricsCategory)
			fTag := e.Type().Field(i).Tag
			if fValue == nil {
				return
			}

			// 请求 Prometheus 查询
			prometheusQueries := fValue.GenerateQuery()
			if prometheusQueries == nil {
				global.DYCLOUD_LOG.Error("prometheusQueries nill failed", zap.Any("err", err))
				return
			}

			//Prometheus健康检查
			prometheusUrl, httpClient, err := m.Health(int(fValue.ClusterId))
			if fValue.Start != 0 || fValue.End != 0 {
				start = fValue.Start / 1000
				end = fValue.End / 1000
			}

			if err != nil {
				global.DYCLOUD_LOG.Error("Prometheus health failed", zap.Any("err", err))
				return
			}

			//http Get请求 Prometheus接口
			promql := url.QueryEscape(prometheusQueries.GetValueByField(fName))
			global.DYCLOUD_LOG.Info(fmt.Sprintf("promql: %s ", prometheusQueries.GetValueByField(fName)))
			fullpromql := fmt.Sprintf("%s/api/v1/query_range?query=%s&start=%d&end=%d&step=%d", prometheusUrl, promql, start, end, step)
			fmt.Println(fullpromql)
			resp, err := httpClient.Get(fullpromql)
			if err != nil {
				global.DYCLOUD_LOG.Error("request metrics data failed", zap.Any("err", err))
				return
			}

			//Prometheus 接口返回数据处理
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				global.DYCLOUD_LOG.Error("read response body failed", zap.Any("err", err))
				return
			}

			var data metrics.PrometheusQueryResp
			if err := json.Unmarshal(body, &data); err != nil {
				global.DYCLOUD_LOG.Error("unmarshal response body to models.PrometheusQueryResp failed", zap.Any("err", err))
				return
			}

			// 配置当前查询的数据结果
			tag := fTag.Get("json")
			tracker.Set(tag[:strings.Index(tag, ",omitempty")], &data)
		}(i)
	}

	// 等待所有查询完成
	wg.Wait()
	return tracker.Metrics, err
}
