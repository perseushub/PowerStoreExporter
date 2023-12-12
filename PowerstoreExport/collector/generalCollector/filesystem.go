package generalCollector

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"powerstore/collector/client"
)

var metricFileCollector = []string{
	"size_total",
	"size_used",
}

// file description
var metricFileDescMap = map[string]string{
	"size_total": "filesystem total size",
	"size_used":  "filesystem used size",
}

type fileCollector struct {
	client  *client.Client
	metrics map[string]*prometheus.Desc
	logger  log.Logger
}

func NewFileCollector(api *client.Client, logger log.Logger) *fileCollector {
	metrics := getFileMetrics(api.IP)
	return &fileCollector{
		client:  api,
		metrics: metrics,
		logger:  logger,
	}
}

func (c *fileCollector) Collect(ch chan<- prometheus.Metric) {
	fileData, err := c.client.GetFile()
	if err != nil {
		level.Warn(c.logger).Log("msg", "get file data error", "err", err)
	}
	for _, file := range gjson.Parse(fileData).Array() {
		name := file.Get("name").String()
		for _, metric := range metricFileCollector {
			value := file.Get(metric)
			metricDesc := c.metrics[metric]
			if value.Exists() && value.Type != gjson.Null {
				ch <- prometheus.MustNewConstMetric(metricDesc, prometheus.GaugeValue, value.Float(), name)
			}
		}
	}
}

func (c *fileCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, descMap := range c.metrics {
		ch <- descMap
	}
}

func getFileMetrics(ip string) map[string]*prometheus.Desc {
	res := map[string]*prometheus.Desc{}
	for _, metric := range metricFileCollector {
		res[metric] = prometheus.NewDesc(
			"powerstore_filesystem_"+metric,
			getFileDescByType(metric),
			[]string{
				"name",
			},
			prometheus.Labels{"IP": ip})
	}
	return res
}

func getFileDescByType(key string) string {
	if v, ok := metricFileDescMap[key]; ok {
		return v
	} else {
		return key
	}
}
