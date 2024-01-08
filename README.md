# README

#### About

One PowerStore Exporter collect multiple PowerStore Devices with PowerStore's Restful API, provide multiple http targets; Support for CNCF Monitoring software Prometheus and Zabbix,they can use these targets, grape PowerStore detail metrices.
This Project test with PowerStore RestAPI version 1.0, 2.0, 3.5; Zabbix version 6.0LTS, Prometheus2.39.1, Grafana8.5.2.

#### Download Compile Run
When you download this project build under golang environment.

```
cd PowerStoreExporter
go build
```
Change  Expoter config file in ./config.yml, you can change this exporter default port 9010  to other port in your local system.
Firstly we strong suggest you to crate operator role user account in PowerStore, then update storeageList section for IP address and PowerStore username/password.


#### Collect metrics and related paths
base path: http://{#PowerstoreExportIP}:{#PowerstoreExportPort}/metrics

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

#### Execution method

```
cd powerstoreExport
#Set the port number,log directory,device IP,username and password in the config.yml
vi config.yml
go build
./powerstoreExport
```

#### 
