package generalCollector

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"powerstore/collector/client"
)

var perfCollectorMetric = []string{
	"avg_read_latency",
	"avg_read_size",
	"avg_latency",
	"avg_write_latency",
	"avg_write_size",
	"read_iops",
	"read_bandwidth",
	"total_iops",
	"total_bandwidth",
	"write_iops",
	"write_bandwidth",
	"avg_io_size",
	"io_workload_cpu_utilization",
}

// performance description
var metricPerfDescMap = map[string]string{
	"avg_read_latency":            "avg latency of read , unit is ms",
	"avg_read_size":               "avg number of read , unit is bytes",
	"avg_latency":                 "avg latency , unit is ms",
	"avg_write_latency":           "avg latency of write , unit is ms",
	"avg_write_size":              "avg number of write , unit is bytes",
	"read_iops":                   "iops of read , unit is iops",
	"read_bandwidth":              "throughput of read , unit is bps",
	"total_iops":                  "iops total , unit is iops",
	"total_bandwidth":             "total throughput , unit is bps",
	"write_iops":                  "iops of write , unit is iops",
	"write_bandwidth":             "throughput of write , unit is bps",
	"avg_io_size":                 "avg number of read and write , unit is bytes",
	"io_workload_cpu_utilization": "usage of CPU for IO workload ",
}

type performanceCollector struct {
	client  *client.Client
	metrics map[string]*prometheus.Desc
	logger  log.Logger
}

func NewPerfCollector(api *client.Client, logger log.Logger) *performanceCollector {
	metrics := getPerfMetrics(api.IP)
	return &performanceCollector{
		client:  api,
		metrics: metrics,
		logger:  logger,
	}
}

func (c *performanceCollector) Collect(ch chan<- prometheus.Metric) {
	idData := client.PowerstoreId[c.client.IP]
	applianceId := idData["appliance"]
	idArray := gjson.Parse(applianceId).Array()
	for _, id := range idArray {
		perfData, err := c.client.GetPerf(id.String())
		if err != nil {
			level.Warn(c.logger).Log("msg", "get perf data error", "err", err)
		}
		perfDataJson := gjson.Parse(perfData)
		perfArray := perfDataJson.Array()
		perf := perfArray[len(perfArray)-1]
		name := perf.Get("appliance_id").String()
		for _, metric := range perfCollectorMetric {
			v := perf.Get(metric)
			metricDesc := c.metrics[metric]
			if v.Exists() && v.Type != gjson.Null {
				ch <- prometheus.MustNewConstMetric(metricDesc, prometheus.GaugeValue, v.Float(), name)
			}
		}
	}

}

func (c *performanceCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, descMap := range c.metrics {
		ch <- descMap
	}
}

func getPerfMetrics(ip string) map[string]*prometheus.Desc {
	res := map[string]*prometheus.Desc{}
	for _, metric := range perfCollectorMetric {
		res[metric] = prometheus.NewDesc(
			"powerstore_perf_"+metric,
			getPerfDescByType(metric),
			[]string{
				"appliance_id",
			},
			prometheus.Labels{"IP": ip})
	}
	return res
}

func getPerfDescByType(key string) string {
	if v, ok := metricPerfDescMap[key]; ok {
		return v
	} else {
		return key
	}
}
