package generalCollector

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"powerstore/collector/client"
)

var metricAppliancePerfCollectorMetric = []string{
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
var metricAppliancePerfDescMap = map[string]string{
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

type metricApplianceCollector struct {
	client  *client.Client
	metrics map[string]*prometheus.Desc
	logger  log.Logger
}

func NewMetricApplianceCollector(api *client.Client, logger log.Logger) *metricApplianceCollector {
	metrics := getMetricApplianceMetrics(api.IP)
	return &metricApplianceCollector{
		client:  api,
		metrics: metrics,
		logger:  logger,
	}
}

func (c *metricApplianceCollector) Collect(ch chan<- prometheus.Metric) {
	applianceArray := client.PowerstoreModuleID[c.client.IP]
	for _, applianceID := range gjson.Parse(applianceArray["appliance"]).Array() {
		id := applianceID.Get("id").String()
		perfData, err := c.client.GetPerf(id)
		if err != nil {
			level.Warn(c.logger).Log("msg", "get appliance performance data error", "err", err)
			continue
		}
		appliancePerformanceArray := gjson.Parse(perfData).Array()
		appliancePerformance := appliancePerformanceArray[len(appliancePerformanceArray)-1]
		name := appliancePerformance.Get("appliance_id").String()
		for _, metricName := range metricAppliancePerfCollectorMetric {
			metricValue := appliancePerformance.Get(metricName)
			metricDesc := c.metrics[metricName]
			if metricValue.Exists() && metricValue.Type != gjson.Null {
				ch <- prometheus.MustNewConstMetric(metricDesc, prometheus.GaugeValue, metricValue.Float(), name)
			}
		}
	}

}

func (c *metricApplianceCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, descMap := range c.metrics {
		ch <- descMap
	}
}

func getMetricApplianceMetrics(ip string) map[string]*prometheus.Desc {
	res := map[string]*prometheus.Desc{}
	for _, metricName := range metricAppliancePerfCollectorMetric {
		res[metricName] = prometheus.NewDesc(
			"powerstore_perf_"+metricName,
			getMetricApplianceDescByType(metricName),
			[]string{
				"appliance_id",
			},
			prometheus.Labels{"IP": ip})
	}
	return res
}

func getMetricApplianceDescByType(key string) string {
	if v, ok := metricAppliancePerfDescMap[key]; ok {
		return v
	} else {
		return key
	}
}
