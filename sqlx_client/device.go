package sqlx_client

import (
	"encoding/json"
	"errors"
	g "github.com/GramYang/gylog"
	"strconv"
	"time"
)

type DeviceDetail struct {
	Id              int    `json:"id" db:"id"`
	DeviceSn        string `json:"deviceSn" db:"device_sn"`
	DeviceId        string `json:"deviceId" db:"device_id"`
	DeviceName      string `json:"deviceName" db:"device_name"`
	DeviceModel     string `json:"deviceModel" db:"device_model"`
	ManufacturerId  string `json:"manufacturerId" db:"manufacturer_id"`
	TenantId        string `json:"tenantId" db:"tenant_id"`
	ProductId       int    `json:"productId" db:"product_id"`
	Imei            string `json:"imei" db:"imei"`
	Imsi            string `json:"imsi" db:"imsi"`
	FirmwareVersion string `json:"firmwareVersion" db:"firmware_version"`
	DeviceVersion   string `json:"deviceVersion" db:"device_version"`
	DeviceStatus    int    `json:"deviceStatus" db:"device_status"`
	AutoObserver    int    `json:"autoObserver" db:"auto_observer"`
	CreateTime      int64  `json:"createTime" db:"create_time"`
	CreateBy        string `json:"createBy" db:"create_by"`
	UpdateTime      int64  `json:"updateTime" db:"update_time"`
	UpdateBy        string `json:"updateBy" db:"update_by"`
	ActiveTime      int64  `json:"activeTime" db:"active_time"`
	LogoutTime      int64  `json:"logoutTime" db:"logout_time"`
	OnlineAt        int64  `json:"onlineAt" db:"online_at"`
	OfflineAt       int64  `json:"offlineAt" db:"offline_at"`
	ProductProtocol int    `json:"productProtocol" db:"product_protocol"`
	SecretKey       string `json:"secretKey" db:"secret_key"`
}

type CreateDeviceDetail struct {
	Code   int          `json:"code"`
	Msg    string       `json:"msg"`
	Result DeviceDetail `json:"result"`
}

func SaveDevice(pid, did string) error {
	if pid == "" || did == "" {
		return errors.New("pid or did null")
	}
	_, err := db.Exec("insert into device(productId,deviceId) values(?,?)", pid, did)
	return err
}

func SaveDeviceDetail(data []byte, para ...string) error {
	mapRes := map[string]interface{}{}
	err := json.Unmarshal(data, &mapRes)
	if err != nil {
		return err
	}
	mapRes1 := mapRes["result"].(map[string]interface{})
	now := time.Now().Unix()
	//ctwing返回的是毫秒数，需要将其清洗成秒数，如果没有创建时间就补上
	if mapRes1["createTime"] != nil {
		mapRes1["createTime"] = int64(mapRes1["createTime"].(float64) / 1000)
	} else {
		mapRes1["createTime"] = now
	}
	if mapRes1["updateTime"] != nil {
		mapRes1["updateTime"] = int64(mapRes1["updateTime"].(float64) / 1000)
	} else {
		mapRes1["updateTime"] = now
	}
	if mapRes1["activeTime"] != nil {
		mapRes1["activeTime"] = int64(mapRes1["activeTime"].(float64) / 1000)
	} else {
		mapRes1["activeTime"] = now
	}
	if mapRes1["logoutTime"] != nil {
		mapRes1["logoutTime"] = int64(mapRes1["logoutTime"].(float64) / 1000)
	} else {
		mapRes1["logoutTime"] = now
	}
	if mapRes1["onlineAt"] != nil {
		mapRes1["onlineAt"] = int64(mapRes1["onlineAt"].(float64) / 1000)
	} else {
		mapRes1["onlineAt"] = now
	}
	if mapRes1["offlineAt"] != nil {
		mapRes1["offlineAt"] = int64(mapRes1["offlineAt"].(float64) / 1000)
	} else {
		mapRes1["offlineAt"] = now
	}
	if len(para) > 0 {
		mapRes1["secretKey"] = para[0]
	}
	data1, err := json.Marshal(&mapRes1)
	if err != nil {
		return err
	}
	var dd DeviceDetail
	err = json.Unmarshal(data1, &dd)
	if err != nil {
		return err
	}
	//time.Time域赋值，因为go不能直接把秒数与time.Time转换，但是time.Time可以和mysql的timestamp互换
	g.Debugln("DeviceDetail: ", dd)
	_, err = db.NamedExec("insert into device_detail(device_sn,device_id,device_name,device_model,"+
		"manufacturer_id,tenant_id,product_id,imei,imsi,firmware_version,device_version,device_status,"+
		"auto_observer,create_time,create_by,update_time,update_by,active_time,logout_time,online_at,"+
		"offline_at,product_protocol,secret_key) "+
		"values(:device_sn,:device_id,:device_name,:device_model,:manufacturer_id,:tenant_id,:product_id,"+
		":imei,:imsi,:firmware_version,:device_version,:device_status,:auto_observer,"+
		"from_unixtime(:create_time),:create_by,from_unixtime(:update_time),:update_by,"+
		"from_unixtime(:active_time),from_unixtime(:logout_time),from_unixtime(:online_at),"+
		"from_unixtime(:offline_at),:product_protocol,:secret_key)", &dd)
	return err
}

//根据设备名称查询，如果置空则返回全部设备，分页。这里就不返回密钥了，密钥专门用一个接口返回。
func QueryDeviceByName(name, pageNow string) ([]DeviceDetail, error) {
	var dds = []DeviceDetail{}
	p, _ := strconv.Atoi(pageNow)
	if p == 0 {
		p = 1
	}
	if name != "" {
		var dd DeviceDetail
		err := db.Get(&dd, "select device_sn,device_id,device_name,device_model,manufacturer_id,"+
			"tenant_id,cast(product_id as signed) as product_id,imei,imsi,firmware_version,device_version,device_status,auto_observer,"+
			"unix_timestamp(create_time) as create_time,create_by,unix_timestamp(update_time) as update_time,"+
			"update_by,unix_timestamp(active_time) as active_time,unix_timestamp(logout_time) as logout_time,"+
			"unix_timestamp(online_at) as online_at,unix_timestamp(offline_at) as offline_at,product_protocol "+
			"from device_detail where device_name=?", name, (p-1)*10, p*10)
		dds = append(dds, dd)
		return dds, err
	} else {
		err := db.Select(&dds, "select device_sn,device_id,device_name,device_model,manufacturer_id,"+
			"tenant_id,cast(product_id as signed) as product_id,imei,imsi,firmware_version,device_version,device_status,auto_observer,"+
			"unix_timestamp(create_time) as create_time,create_by,unix_timestamp(update_time) as update_time,"+
			"update_by,unix_timestamp(active_time) as active_time,unix_timestamp(logout_time) as logout_time,"+
			"unix_timestamp(online_at) as online_at,unix_timestamp(offline_at) as offline_at,product_protocol "+
			"from device_detail order by id desc limit ?,?", (p-1)*10, p*10)
		return dds, err
	}
}
