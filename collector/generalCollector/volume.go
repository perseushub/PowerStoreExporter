package generalCollector

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"powerstore/collector/client"
)

var volumeCollectorMetrics = []string{
	"state",
	"size",
	"logical_used",
}

var metricVolumeDescMap = map[string]string{
	"state":        "1 is ready ,0 is other",
	"size":         "the unit is B",
	"logical_used": "the unit is B",
}

var statuVolumeMetricsMap = map[string]map[string]int{
	"state": {"Ready": 1, "other": 0},
}

type volumeCollector struct {
	client  *client.Client
	metrics map[string]*prometheus.Desc
	logger  log.Logger
}

func NewVolumeCollector(api *client.Client, logger log.Logger) *volumeCollector {
	metrics := getVolumeMetrics(api.IP)
	return &volumeCollector{
		client:  api,
		metrics: metrics,
		logger:  logger,
	}
}

func (c *volumeCollector) Collect(ch chan<- prometheus.Metric) {
	volumeData, err := c.client.GetVolume(c.client.Version)
	if err != nil {
		level.Warn(c.logger).Log("msg", "get volume data error", "err", err)
	}
	volumeDataJson := gjson.Parse(volumeData)
	volumeArray := volumeDataJson.Array()
	for _, volume := range volumeArray {
		name := volume.Get("name").String()
		for _, metric := range volumeCollectorMetrics {
			rs := volume.Get(metric)
			value := getVolumeFloatDate(metric, rs)
			metricDesc := c.metrics[metric]
			if rs.Exists() && rs.Type != gjson.Null {
				ch <- prometheus.MustNewConstMetric(metricDesc, prometheus.GaugeValue, value, name)
			}
		}
	}
}

func (c *volumeCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, descMap := range c.metrics {
		ch <- descMap
	}
}

func getVolumeFloatDate(key string, value gjson.Result) float64 {
	if v, ok := statuVolumeMetricsMap[key]; ok {
		if res, ok2 := v[value.String()]; ok2 {
			return float64(res)
		} else {
			return float64(v["other"])
		}
	} else {
		return value.Float()
	}
}

func getVolumeMetrics(ip string) map[string]*prometheus.Desc {
	res := map[string]*prometheus.Desc{}
	for _, metric := range volumeCollectorMetrics {
		res[metric] = prometheus.NewDesc(
			"powerstore_volume_"+metric,
			getVolumeDescByType(metric),
			[]string{"name"},
			prometheus.Labels{"IP": ip})
	}

	return res
}

func getVolumeDescByType(key string) string {
	if v, ok := metricVolumeDescMap[key]; ok {
		return v
	} else {
		return key
	}
}
