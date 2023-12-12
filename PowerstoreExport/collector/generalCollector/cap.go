package generalCollector

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"powerstore/collector/client"
)

var capCollectorMetric = []string{
	"last_logical_provisioned",
	"last_logical_used",
	"last_physical_total",
	"last_physical_used",
	"max_logical_provisioned",
	"max_logical_used",
	"max_physical_total",
	"max_physical_used",
	"last_data_physical_used",
	"max_data_physical_used",
	"last_efficiency_ratio",
	"last_data_reduction",
	"last_snapshot_savings",
	"last_thin_savings",
	"max_efficiency_ratio",
	"max_data_reduction",
	"max_snapshot_savings",
	"max_thin_savings",
	"last_shared_logical_used",
	"max_shared_logical_used",
}

var metricCapDescMap = map[string]string{
	"last_logical_provisioned": "last logical provisioned,unit is B",
	"last_logical_used":        "last logical has been used,unit is B",
	"last_physical_total":      "total last physical ,unit is B",
	"last_physical_used":       "last physical has been used ,unit is B",
	"max_logical_provisioned":  "max logical provisioned,unit is B",
	"max_logical_used":         "max used logical,unit is B",
	"max_physical_total":       "max total physical ,unit is B",
	"max_physical_used":        "max used physical,unit is B",
	"last_data_physical_used":  "last data used physical,unit is B",
	"max_data_physical_used":   "max used data physical used,unit is B",
	"last_efficiency_ratio":    "last efficiency ratio,:1",
	"last_data_reduction":      "last data reduction",
	"last_snapshot_savings":    "last snapshot savings",
	"last_thin_savings":        "last thin savings",
	"max_efficiency_ratio":     "max efficiency ratio :1",
	"max_data_reduction":       "max data reduction,unit is B",
	"max_snapshot_savings":     "max snapshot savings",
	"max_thin_savings":         "max thin savings",
	"last_shared_logical_used": "last shared logical used,unit is B",
	"max_shared_logical_used":  "max shared logical used,unit is B",
}

type capCollector struct {
	client  *client.Client
	metrics map[string]*prometheus.Desc
	logger  log.Logger
}

func NewCapCollector(api *client.Client, logger log.Logger) *capCollector {
	metrics := getCapMetrics(api.IP)
	return &capCollector{
		client:  api,
		metrics: metrics,
		logger:  logger,
	}
}

func (c *capCollector) Collect(ch chan<- prometheus.Metric) {
	idData := client.PowerstoreId[c.client.IP]
	applianceId := idData["appliance"]
	idArray := gjson.Parse(applianceId).Array()
	for _, id := range idArray {
		capData, err := c.client.GetCap(id.String())
		if err != nil {
			level.Warn(c.logger).Log("msg", "get cap data error", "err", err)
		}
		capDataJson := gjson.Parse(capData)
		capDataArray := capDataJson.Array()
		capa := capDataArray[len(capDataArray)-1]
		name := capa.Get("appliance_id").String()
		for _, metric := range capCollectorMetric {
			v := capa.Get(metric)
			metricDesc := c.metrics[metric]
			if v.Exists() && v.Type != gjson.Null {
				ch <- prometheus.MustNewConstMetric(metricDesc, prometheus.GaugeValue, v.Float(), name)
			}
		}
	}

}

func (c *capCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, descMap := range c.metrics {
		ch <- descMap
	}
}

func getCapMetrics(ip string) map[string]*prometheus.Desc {
	res := map[string]*prometheus.Desc{}
	for _, metric := range capCollectorMetric {
		res[metric] = prometheus.NewDesc(
			"powerstore_cap_"+metric,
			getCapDescByType(metric),
			[]string{
				"appliance_id",
			},
			prometheus.Labels{"IP": ip})
	}
	return res
}

func getCapDescByType(key string) string {
	if v, ok := metricCapDescMap[key]; ok {
		return v
	} else {
		return key
	}
}
