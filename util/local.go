package util

import (
	"io/ioutil"
	"os"
)

const NEWDEVICES="template/new_devices.txt"

//在数据库出错的情况下将新建设备的数据备份到本地
func WriteDeviceToLocal(pid,did string){
	if !isExist(NEWDEVICES){
		_=os.MkdirAll(NEWDEVICES,os.ModePerm)
	}
	_=ioutil.WriteFile(NEWDEVICES,[]byte(pid+":"+did+"\n"),0666)
}

func isExist(path string)bool{
	_,err:=os.Stat(path)
	if err!=nil{
		if os.IsExist(err){
			return true
		}
		return false
	}
	return true
}