package generalCollector

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"powerstore/collector/client"
)

type metricWearMetricCollector struct {
	client  *client.Client
	metrics map[string]*prometheus.Desc
	logger  log.Logger
}

func NewWearMetricCollector(api *client.Client, logger log.Logger) *metricWearMetricCollector {
	metrics := getWearMetrics(api.IP)
	return &metricWearMetricCollector{
		client:  api,
		metrics: metrics,
		logger:  logger,
	}
}

func (c *metricWearMetricCollector) Collect(ch chan<- prometheus.Metric) {
	idData := client.PowerstoreId[c.client.IP]
	driveId := idData["drive"]
	idArray := gjson.Parse(driveId).Array()
	for _, drId := range idArray {
		id := drId.Get("id").String()
		name := drId.Get("name").String()
		wearMetricData, err := c.client.GetWearMetricByDrive(id)
		if err != nil {
			level.Warn(c.logger).Log("msg", "get metric Wear data error", "err", err)
			continue
		}
		metricWearArray := gjson.Parse(wearMetricData).Array()
		arrLen := len(metricWearArray)
		if arrLen == 0 {
			continue
		}
		wearData := metricWearArray[arrLen-1]
		v := wearData.Get("percent_endurance_remaining")
		metricDesc := c.metrics["wear"]
		if v.Exists() && v.Type != gjson.Null {
			ch <- prometheus.MustNewConstMetric(metricDesc, prometheus.GaugeValue, v.Float(), name)
		}
	}

}

func (c *metricWearMetricCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, descMap := range c.metrics {
		ch <- descMap
	}
}

func getWearMetrics(ip string) map[string]*prometheus.Desc {
	res := map[string]*prometheus.Desc{}

	res["wear"] = prometheus.NewDesc(
		"powerstore_wear_metrics_by_drive",
		"this is the percent of endurance remaining about drives",
		[]string{
			"name",
		},
		prometheus.Labels{"IP": ip})
	return res
}
