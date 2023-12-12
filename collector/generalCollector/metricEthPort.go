package generalCollector

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"powerstore/collector/client"
)

var metricEthPortCollectorMetric = []string{
	"bytes_rx_ps",
	"bytes_tx_ps",
	"pkt_rx_crc_error_ps",
	"pkt_rx_no_buffer_error_ps",
	"pkt_rx_ps",
	"pkt_tx_error_ps",
	"pkt_tx_ps",
}

var metricMetricEthPortDescMap = map[string]string{
	"bytes_rx_ps":               "receive bytes in a second",
	"bytes_tx_ps":               "send bytes in a second",
	"pkt_rx_crc_error_ps":       "packet receive crc error in a second",
	"pkt_rx_no_buffer_error_ps": "packet receive no buffer error in a second",
	"pkt_rx_ps":                 "packet receive in a second",
	"pkt_tx_error_ps":           "packet send error in a second",
	"pkt_tx_ps":                 "packet get in a second",
}

type metricEthPortCollector struct {
	client  *client.Client
	metrics map[string]*prometheus.Desc
	logger  log.Logger
}

func NewMetricEthPortCollector(api *client.Client, logger log.Logger) *metricEthPortCollector {
	metrics := getMetricEthPortfMetrics(api.IP)
	return &metricEthPortCollector{
		client:  api,
		metrics: metrics,
		logger:  logger,
	}
}

func (c *metricEthPortCollector) Collect(ch chan<- prometheus.Metric) {
	idData := client.PowerstoreId[c.client.IP]
	ethPortId := idData["ethport"]
	idArray := gjson.Parse(ethPortId).Array()
	for _, portId := range idArray {
		id := portId.Get("id").String()
		name := portId.Get("name").String()
		ethPortData, err := c.client.GetMetricEthPort(id)
		if err != nil {
			level.Warn(c.logger).Log("msg", "get metric ethPort data error", "err", err)
		}
		ethPortArray := gjson.Parse(ethPortData).Array()
		arrLen := len(ethPortArray)
		if arrLen == 0 {
			continue
		}
		ethport := ethPortArray[arrLen-1]
		for _, metric := range metricEthPortCollectorMetric {
			v := ethport.Get(metric)
			metricDesc := c.metrics[metric]
			if v.Exists() && v.Type != gjson.Null {
				ch <- prometheus.MustNewConstMetric(metricDesc, prometheus.GaugeValue, v.Float(), name)
			}
		}
	}
}

func (c *metricEthPortCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, descMap := range c.metrics {
		ch <- descMap
	}
}

func getMetricEthPortfMetrics(ip string) map[string]*prometheus.Desc {
	res := map[string]*prometheus.Desc{}
	for _, metric := range metricEthPortCollectorMetric {
		res[metric] = prometheus.NewDesc(
			"powerstore_metricEthPort_"+metric,
			getMetricEthPortDescByType(metric),
			[]string{
				"eth_port_id",
			},
			prometheus.Labels{"IP": ip})
	}
	return res
}

func getMetricEthPortDescByType(key string) string {
	if v, ok := metricMetricEthPortDescMap[key]; ok {
		return v
	} else {
		return key
	}
}
