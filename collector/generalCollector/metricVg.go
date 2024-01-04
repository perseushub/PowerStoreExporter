package generalCollector

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"powerstore/collector/client"
)

var metricVgCollectorMetric = []string{
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
}

var metricMetricVgDescMap = map[string]string{
	"avg_read_latency":  "avg latency time of read,unit is ms",
	"avg_read_size":     "avg size of read,unit is B",
	"avg_latency":       "avg latency time,unit is ms",
	"avg_write_latency": "avg latency time of write,unit is ms",
	"avg_write_size":    "avg size of write,unit is B",
	"read_iops":         "iops of read,unit is iops",
	"read_bandwidth":    "bandwidth of read,unit is bps",
	"total_iops":        "total iops,unit is iops",
	"total_bandwidth":   "total bandwidth,unit is bps",
	"write_iops":        "iops of write,unit is iops",
	"write_bandwidth":   "bandwidth of write,unit is bps",
	"avg_io_size":       "avg size of IO",
}

type metricVgCollector struct {
	client  *client.Client
	metrics map[string]*prometheus.Desc
	logger  log.Logger
}

func NewMetricVgCollector(api *client.Client, logger log.Logger) *metricVgCollector {
	metrics := getMetricVgfMetrics(api.IP)
	return &metricVgCollector{
		client:  api,
		metrics: metrics,
		logger:  logger,
	}
}

func (c *metricVgCollector) Collect(ch chan<- prometheus.Metric) {
	vgArray := client.PowerstoreModuleID[c.client.IP]
	for _, vgId := range gjson.Parse(vgArray["volumegroup"]).Array() {
		id := vgId.Get("id").String()
		name := vgId.Get("name").String()
		metricVgData, err := c.client.GetMetricVg(id)
		if err != nil {
			level.Warn(c.logger).Log("msg", "get volume group performance data error", "err", err)
			continue
		}
		vgDataArray := gjson.Parse(metricVgData).Array()
		if len(vgDataArray) == 0 {
			continue
		}
		vgData := vgDataArray[len(vgDataArray)-1]
		for _, metricName := range metricVgCollectorMetric {
			metricValue := vgData.Get(metricName)
			metricDesc := c.metrics[metricName]
			if metricValue.Exists() && metricValue.Type != gjson.Null {
				ch <- prometheus.MustNewConstMetric(metricDesc, prometheus.GaugeValue, metricValue.Float(), name)
			}
		}
	}
}

func (c *metricVgCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, descMap := range c.metrics {
		ch <- descMap
	}
}

func getMetricVgfMetrics(ip string) map[string]*prometheus.Desc {
	res := map[string]*prometheus.Desc{}
	for _, metricName := range metricVgCollectorMetric {
		res[metricName] = prometheus.NewDesc(
			"powerstore_metricVg_"+metricName,
			getMetricVgDescByType(metricName),
			[]string{
				"volume_group_id",
			},
			prometheus.Labels{"IP": ip})
	}
	return res
}

func getMetricVgDescByType(key string) string {
	if v, ok := metricMetricVgDescMap[key]; ok {
		return v
	} else {
		return key
	}
}
