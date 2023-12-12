package client

import (
	"encoding/json"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"powerstore/utils"
)

var PowerstoreId = map[string]map[string]string{}

var ParamMap = make(map[string]interface{})

func (c *Client) getData(path, method, body string) (string, error) {
	utils.ReqCounter <- 1
	result, err := c.getResource(method, path, body)
	<-utils.ReqCounter
	return result, err
}

func (c *Client) GetCluster() (string, error) {
	return c.getData("cluster?select=*", "GET", "")
}

func (c *Client) GetPort(name string) (string, error) {
	return c.getData(name+"?select=*", "GET", "")
}

func (c *Client) GetHardware(id string) (string, error) {
	return c.getData("hardware?select=*&type=eq."+id, "GET", "")
}

func (c *Client) GetVolume(version string) (string, error) {
	if version == "v3" {
		return c.getData("volume_list_cma_view?select=*", "GET", "")
	}
	return c.getData("volume?select=*", "GET", "")
}

func (c *Client) GetAppliance() (string, error) {
	return c.getData("appliance?select=*", "GET", "")
}

func (c *Client) GetFile() (string, error) {
	return c.getData("file_system?select=*", "GET", "")
}

func (c *Client) GetNas() (string, error) {
	return c.getData("nas_server?select=*", "GET", "")
}

func (c *Client) GetVolumeGroup() (string, error) {
	return c.getData("volume_group_list_cma_view?select=*", "GET", "")
}

func (c *Client) GetPerf(id string) (string, error) {
	param := ParamMap
	param["entity"] = "performance_metrics_by_appliance"
	param["entity_id"] = id
	valueBody, _ := json.Marshal(param)
	return c.getData("metrics/generate", "POST", string(valueBody))
}

func (c *Client) GetCap(id string) (string, error) {
	param := ParamMap
	param["entity"] = "space_metrics_by_appliance"
	param["entity_id"] = id
	param["interval"] = "One_Day"
	valueBody, _ := json.Marshal(param)
	return c.getData("metrics/generate", "POST", string(valueBody))
}

func (c *Client) GetVg(id string) (string, error) {
	param := ParamMap
	param["entity"] = "performance_metrics_by_vg"
	param["entity_id"] = id
	valueBody, _ := json.Marshal(param)
	return c.getData("metrics/generate", "POST", string(valueBody))
}

func (c *Client) GetMetricVolume(id string) (string, error) {
	param := ParamMap
	param["entity"] = "performance_metrics_by_volume"
	param["entity_id"] = id
	valueBody, _ := json.Marshal(param)
	return c.getData("metrics/generate", "POST", string(valueBody))
}

func (c *Client) GetMetricFcPort(id string) (string, error) {
	param := ParamMap
	param["entity"] = "performance_metrics_by_fe_fc_port"
	param["entity_id"] = id
	valueBody, _ := json.Marshal(param)
	return c.getData("metrics/generate", "POST", string(valueBody))
}

func (c *Client) GetMetricEthPort(id string) (string, error) {
	param := ParamMap
	param["entity"] = "performance_metrics_by_fe_eth_port"
	param["entity_id"] = id
	valueBody, _ := json.Marshal(param)
	return c.getData("metrics/generate", "POST", string(valueBody))
}

func (c *Client) GetMetricAppliance(id string) (string, error) {
	param := ParamMap
	param["entity"] = "performance_metrics_by_appliance"
	param["entity_id"] = id
	valueBody, _ := json.Marshal(param)
	return c.getData("metrics/generate", "POST", string(valueBody))
}

func (c *Client) GetWearMetricByDrive(id string) (string, error) {
	param := ParamMap
	param["entity"] = "wear_metrics_by_drive"
	param["entity_id"] = id
	valueBody, _ := json.Marshal(param)
	return c.getData("metrics/generate", "POST", string(valueBody))
}

func (c *Client) GetApplianceId() (string, error) {
	return c.getData("appliance?select=id,name", "GET", "")
}

func (c *Client) GetVolumeGroupId() (string, error) {
	return c.getData("volume_group?select=id,name", "GET", "")
}

func (c *Client) GetVolumeId(version string) (string, error) {
	if version == "v3" {
		return c.getData("volume_list_cma_view?select=id,name", "GET", "")
	}
	return c.getData("volume?select=id,name", "GET", "")
}

func (c *Client) GetEthPortId() (string, error) {
	return c.getData("eth_port?select=id,name", "GET", "")
}

func (c *Client) GetFcPortId() (string, error) {
	return c.getData("fc_port?select=id,name", "GET", "")
}

func (c *Client) GetDrivesId() (string, error) {
	return c.getData("hardware?select=id,name", "GET", "")
}

func (c *Client) Init(logger log.Logger) {
	id := make(map[string]string)
	ApplianceId, err := c.GetApplianceId()
	if err != nil {
		level.Error(logger).Log("msg", "GetApplianceId error", "err", err, "ip", c.IP)
	}
	id["appliance"] = ApplianceId
	VolumeId, err := c.GetVolumeId(c.Version)
	if err != nil {
		level.Error(logger).Log("msg", "GetVolumeId error", "err", err, "ip", c.IP)
	}
	id["volume"] = VolumeId
	VolumeGroupId, err := c.GetVolumeGroupId()
	if err != nil {
		level.Error(logger).Log("msg", "GetVolumeGroupId error", "err", err, "ip", c.IP)
	}
	id["volumegroup"] = VolumeGroupId
	EthPortId, err := c.GetEthPortId()
	if err != nil {
		level.Error(logger).Log("msg", "GetEthPortId error", "err", err, "ip", c.IP)
	}
	id["ethport"] = EthPortId
	FcPortId, err := c.GetFcPortId()
	if err != nil {
		level.Error(logger).Log("msg", "GetFcPortId error", "err", err, "ip", c.IP)
	}
	id["fcport"] = FcPortId
	DrivesId, err := c.GetDrivesId()
	if err != nil {
		level.Error(logger).Log("msg", "GetDrivesId error", "err", err, "ip", c.IP)
	}
	id["drive"] = DrivesId
	PowerstoreId[c.IP] = id
}
