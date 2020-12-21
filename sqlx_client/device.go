package sqlx_client

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type DeviceDetail struct {
	Id int `json:"id" db:"id"`
	DeviceSn string `json:"deviceSn" db:"device_sn"`
	DeviceId string `json:"deviceId" db:"device_id"`
	DeviceName string `json:"deviceName" db:"device_name"`
	DeviceModel string `json:"deviceModel" db:"device_model"`
	ManufacturerId string `json:"manufacturerId" db:"manufacturer_id"`
	TenantId string `json:"tenantId" db:"tenant_id"`
	ProductId string `json:"productId" db:"product_id"`
	Imei string `json:"imei" db:"imei"`
	Imsi string `json:"imsi" db:"imsi"`
	FirmwareVersion string `json:"firmwareVersion" db:"firmware_version"`
	DeviceVersion string `json:"deviceVersion" db:"device_version"`
	DeviceStatus int `json:"deviceStatus" db:"device_status"`
	AutoObserver int `json:"autoObserver" db:"auto_observer"`
	CreateTime time.Time `json:"createTime" db:"create_time"`
	CreateBy string `json:"createBy" db:"create_by"`
	UpdateTime time.Time `json:"updateTime" db:"update_time"`
	UpdateBy string `json:"updateBy" db:"update_by"`
	ActiveTime time.Time `json:"activeTime" db:"active_time"`
	LogoutTime time.Time `json:"logoutTime" db:"logout_time"`
	OnlineAt time.Time `json:"onlineAt" db:"online_at"`
	OfflineAt time.Time `json:"offlineAt" db:"offline_at"`
	ProductProtocol int `json:"productProtocol" db:"product_protocol"`
}

func SaveDevice(pid,did string) error{
	if pid==""||did==""{
		return errors.New("pid or did null")
	}
	_,err:=db.Exec("insert into device(productId,deviceId) values(?,?)",pid,did)
	return err
}

func SaveDeviceDetail(data []byte) error{
	var dd DeviceDetail
	err:=json.Unmarshal(data,&dd)
	if err!=nil{
		return err
	}
	//如果没有创建时间就补上
	if(dd.CreateTime==time.Time{}){
		dd.CreateTime=time.Now()
	}
	_,err=db.NamedExec("insert into device_detail(device_sn,device_id,device_name,device_model," +
		"manufacturer_id,tenant_id,product_id,imei,imsi,firmware_version,device_version,device_status," +
		"auto_observer,create_time,create_by,update_time,update_by,active_time,logout_time,online_at," +
		"offline_at,product_protocol) " +
		"values(:device_sn,:device_id,:device_name,:device_model,:manufacturer_id,:tenant_id,:product_id," +
		":imei,:imsi,:firmware_version,:device_version,:device_status,:auto_observer,:create_time," +
		":create_by,:update_time,:update_by,:active_time,:logout_time,:online_at,:offline_at," +
		":product_protocol)",&dd)
	return err
}

//根据设备名称查询，如果置空则返回全部设备，分页
func QueryDeviceByName(name,pageNow string) ([]DeviceDetail,error){
	var dds =[]DeviceDetail{}
	p,_:=strconv.Atoi(pageNow)
	if p==0{
		p=1
	}
	if name!=""{
		var dd DeviceDetail
		err:=db.Get(&dd,"select * from device_detail where device_name=? order by id desc limit ?,?",
			name,(p-1)*10,p*10)
		dds=append(dds,dd)
		return dds,err
	}else{
		err:=db.Select(&dds,"select * from device_detail order by id desc limit ?,?",
			(p-1)*10,p*10)
		return dds,err
	}
}