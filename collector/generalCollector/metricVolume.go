package generalCollector

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"powerstore/collector/client"
)

var metricVolumeCollectorMetric = []string{
	"avg_read_latency",
	"avg_latency",
	"avg_write_latency",
	"read_iops",
	"read_bandwidth",
	"total_iops",
	"total_bandwidth",
	"write_iops",
	"write_bandwidth",
}

var metricMetricVolumeDescMap = map[string]string{
	"avg_read_latency":  "avg latency time of read,unit is ms",
	"avg_latency":       "avg latency time,unit is ms",
	"avg_write_latency": "avg latency time of write,unit is ms",
	"read_iops":         "iops of read,unit is iops",
	"read_bandwidth":    "bandwidth of read,unit is bps",
	"total_iops":        "total iops,unit is iops",
	"total_bandwidth":   "total bandwidth,unit is bps",
	"write_iops":        "iops of write,unit is iops",
	"write_bandwidth":   "bandwidth of write,unit is bps",
}

type metricVolumeCollector struct {
	client  *client.Client
	metrics map[string]*prometheus.Desc
	logger  log.Logger
}

func NewMetricVolumeCollector(api *client.Client, logger log.Logger) *metricVolumeCollector {
	metrics := getMetricVolumefMetrics(api.IP)
	return &metricVolumeCollector{
		client:  api,
		metrics: metrics,
		logger:  logger,
	}
}

func (c *metricVolumeCollector) Collect(ch chan<- prometheus.Metric) {
	idData := client.PowerstoreId[c.client.IP]
	volumeId := idData["volume"]
	idArray := gjson.Parse(volumeId).Array()
	for _, volId := range idArray {
		id := volId.Get("id").String()
		name := volId.Get("name").String()
		metricVolData, err := c.client.GetMetricVolume(id)
		if err != nil {
			level.Warn(c.logger).Log("msg", "get metric Vg data error", "err", err)
		}
		metricVolArray := gjson.Parse(metricVolData).Array()
		arrLen := len(metricVolArray)
		if arrLen == 0 {
			continue
		}
		volData := metricVolArray[arrLen-1]
		for _, metric := range metricVolumeCollectorMetric {
			v := volData.Get(metric)
			metricDesc := c.metrics[metric]
			if v.Exists() && v.Type != gjson.Null {
				ch <- prometheus.MustNewConstMetric(metricDesc, prometheus.GaugeValue, v.Float(), name)
			}
		}
	}
}

func (c *metricVolumeCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, descMap := range c.metrics {
		ch <- descMap
	}
}

func getMetricVolumefMetrics(ip string) map[string]*prometheus.Desc {
	res := map[string]*prometheus.Desc{}
	for _, metric := range metricVgCollectorMetric {
		res[metric] = prometheus.NewDesc(
			"powerstore_metricVolume_"+metric,
			getMetricVolumeDescByType(metric),
			[]string{
				"volume_id",
			},
			prometheus.Labels{"IP": ip})
	}
	return res
}

func getMetricVolumeDescByType(key string) string {
	if v, ok := metricMetricVolumeDescMap[key]; ok {
		return v
	} else {
		return key
	}
}
