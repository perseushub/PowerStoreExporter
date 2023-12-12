package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"powerstore/collector/client"
	"powerstore/collector/generalCollector"
	"powerstore/utils"
	"strconv"
)

func Run(config *utils.Config, logger log.Logger) {
	r := gin.New()
	r.Use(gin.Recovery())
	for _, storage := range config.StorageList {
		client, err := client.NewClient(storage, logger)
		if err != nil {
			level.Error(logger).Log("msg", "init Powerstore client error", "err", err, "ip", storage.Ip)
		}

		client.Init(logger)

		ClusterRegistry := prometheus.NewPedanticRegistry()
		PortRegistry := prometheus.NewPedanticRegistry()
		FileRegistry := prometheus.NewPedanticRegistry()
		HardwareRegistry := prometheus.NewPedanticRegistry()
		VolumeRegistry := prometheus.NewPedanticRegistry()
		ApplianceRegistry := prometheus.NewPedanticRegistry()
		NasRegistry := prometheus.NewPedanticRegistry()
		VolumeGroupRegistry := prometheus.NewPedanticRegistry()
		CapRegistry := prometheus.NewPedanticRegistry()

		ClusterRegistry.MustRegister(generalCollector.NewClusterCollector(client, logger))
		ClusterRegistry.MustRegister(generalCollector.NewPerfCollector(client, logger))
		PortRegistry.MustRegister(generalCollector.NewPortCollector(client, logger))
		PortRegistry.MustRegister(generalCollector.NewMetricFcPortCollector(client, logger))
		PortRegistry.MustRegister(generalCollector.NewMetricEthPortCollector(client, logger))
		FileRegistry.MustRegister(generalCollector.NewFileCollector(client, logger))
		HardwareRegistry.MustRegister(generalCollector.NewHardwareCollector(client, logger))
		HardwareRegistry.MustRegister(generalCollector.NewWearMetricCollector(client, logger))
		VolumeRegistry.MustRegister(generalCollector.NewVolumeCollector(client, logger))
		VolumeRegistry.MustRegister(generalCollector.NewMetricVolumeCollector(client, logger))
		ApplianceRegistry.MustRegister(generalCollector.NewApplianceCollector(client, logger))
		NasRegistry.MustRegister(generalCollector.NewNasCollector(client, logger))
		VolumeGroupRegistry.MustRegister(generalCollector.NewVolumeGroupCollector(client, logger))
		VolumeGroupRegistry.MustRegister(generalCollector.NewMetricVgCollector(client, logger))
		CapRegistry.MustRegister(generalCollector.NewCapCollector(client, logger))

		metricsGroup := r.Group(fmt.Sprintf("/metrics/%s", storage.Ip))
		{
			metricsGroup.GET("/cluster", utils.PrometheusHandler(ClusterRegistry, logger))
			metricsGroup.GET("/port", utils.PrometheusHandler(PortRegistry, logger))
			metricsGroup.GET("/file", utils.PrometheusHandler(FileRegistry, logger))
			metricsGroup.GET("/hardware", utils.PrometheusHandler(HardwareRegistry, logger))
			metricsGroup.GET("/volume", utils.PrometheusHandler(VolumeRegistry, logger))
			metricsGroup.GET("/appliance", utils.PrometheusHandler(ApplianceRegistry, logger))
			metricsGroup.GET("/nas", utils.PrometheusHandler(NasRegistry, logger))
			metricsGroup.GET("/volumeGroup", utils.PrometheusHandler(VolumeGroupRegistry, logger))
			metricsGroup.GET("/capacity", utils.PrometheusHandler(CapRegistry, logger))

		}
	}

	// exporter Performance
	r.GET("/performance", func(context *gin.Context) {
		h := promhttp.Handler()
		h.ServeHTTP(context.Writer, context.Request)
	})

	httpPort := fmt.Sprintf(":%s", strconv.Itoa(config.Exporter.Port))
	level.Info(logger).Log("msg", "~~~~~~~~~~~~~Start Powerstore Exporter~~~~~~~~~~~~~~")
	r.Run(httpPort)
}
