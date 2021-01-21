package ctwing_client

import (
	aep "apis/Aep_device_management"
	"encoding/json"
	g "github.com/GramYang/gylog"
	"net/http"
)

type device struct {
	DeviceName string `json:"deviceName"`
	DeviceSn   string `json:"deviceSn"`
	Imei       string `json:"imei"`
	Operator   string `json:"operator"`
	Other      other  `json:"other"`
	ProductId  int    `json:"productId"`
}

type other struct {
	AutoObserver int    `json:"autoObserver"`
	Imsi         string `json:"imsi"`
	PskValue     string `json:"pskValue"`
}

func CreateDevice(body []byte) (*http.Response, error) {
	var d device
	_ = json.Unmarshal(body, &d)
	g.Debugln("create device body: ", d)
	data, _ := json.Marshal(&d)
	return aep.CreateDevice(AppKey, AppSecret, MasterKey, string(data))
}

func UpdateDevice(body []byte) (*http.Response, error) {
	var mapBody = make(map[string]interface{})
	err := json.Unmarshal(body, &mapBody)
	if err != nil {
		return nil, err
	}
	deviceId := mapBody["deviceId"].(string)
	delete(mapBody, "deviceId")
	data, _ := json.Marshal(mapBody)
	resp, err := aep.UpdateDevice(AppKey, AppSecret, MasterKey, deviceId, string(data))
	return resp, err
}

func DeleteDevice(pid, dids string) (*http.Response, error) {
	return aep.DeleteDevice(AppKey, AppSecret, MasterKey, pid, dids)
}

func QuerySingleDevice(pid, did string) (*http.Response, error) {
	return aep.QueryDevice(AppKey, AppSecret, MasterKey, did, pid)
}

func QueryDevices(pid, sv, page string) (*http.Response, error) {
	//这里的pageNow，0和1都是一样的
	return aep.QueryDeviceList(AppKey, AppSecret, MasterKey, pid, sv, page, "10")
}
