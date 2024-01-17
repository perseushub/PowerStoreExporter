# README

#### About

PowerStore Exporter for Both Prometheus and Zabbix.

One PowerStore Exporter collect multiple PowerStore Devices with PowerStore's Restful API, provide multiple http targets; Support for CNCF Monitoring software Prometheus and Zabbix,they can use these targets, scrape PowerStore detail metrices.
This Project test with PowerStore RestAPI version 1.0, 2.0, 3.5; Zabbix version 6.0LTS, Prometheus2.39.1, Grafana9.3.8.

#### Download Compile
When you download this project build under golang environment.

```
cd PowerStoreExporter
go build -o PowerStoreExporter
```
#### Run
Change Expoter config file in ./config.yml, you can change this exporter default port 9010  to other port in your local system.
Firstly we strong suggest you to crate operator role user account in PowerStore, then update storeageList section for IP address and PowerStore username/password.
```
./PowerStoreExporter -c config.yml
```


#### Collect metrics and related paths
base path: http://{#PowerStore Exporter IP}:{#PowerStore Exporter Port}/metrics

```
Cluster              /{#PowerStoreIP}/cluster
Appliance            /{#PowerStoreIP}/appliance
Capacity             /{#PowerStoreIP}/capacity
Hardware             /{#PowerStoreIP}/hardware
Volume               /{#PowerStoreIP}/volume
VolumeGroup          /{#PowerStoreIP}/volumeGroup
Port                 /{#PowerStoreIP}/port
Nas                  /{#PowerStoreIP}/nas
FileSystem           /{#PowerStoreIP}/file
```
Sample: http://127.0.0.1:9010/metrics/10.0.0.1/Cluster

You can chose one of Prometheus or Zabbix monitoring software to scrape this exporter targets, then use Grafana to render the metrics.
For Prometheus user: PowerStores --> PowerStore Expoter --> multiple targets --> Prometheus scrape jobs --> Promtheus --> Grafana

For Zabbix user: PowerStores --> PowerStore Expoter --> multiple targets --> [ Create PowerStore host in Zabbix --> Link this host with PowerStore Zabbix template --> Scrape targets by Zabbix http client --> Zabbix DB --> Zabbix API] --> Grafana

#### Prometheus and Grafana
Add ./templates/prometheus/prometheus.yml  all jobs to your Prometheus yml config file, then restart your Prometheus or reload, your can update scrape interval time support application monitoring requirement.
We use Grafana Dashboard to render Prometheus metrics data.

#### Zabbix and Grafana
When you create host in Zabbix system, use ./templates/zabbix/zbx_exporter_tempaltes.yaml to link PowerStore Zabbix host.
You can use Zabbix latest date to get those detail metrics data. 
We use Zabbix as Grafana datasource then render PowerStore metrics data in Grafna Dashboard.
You alse can create Zabbix dashboard then add widget get PowerStore item as you need.

