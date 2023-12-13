# README

#### About

About the development of the code for the collection engine of Powerstore device metrics.

#### Collect metrics and related paths
base path: http://{#PowerstoreExportIP}/metrics

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
