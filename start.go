package main

import (
	"github.com/gin-gonic/gin"
	"iot_admin/config"
	"iot_admin/log"
	"iot_admin/router"
	sc "iot_admin/sqlx_client"
)



func main(){
	log.InitLog2(1)
	config.SetUp()
	sc.SetUp()
	r:=router.NewRouter()
	gin.SetMode(config.Conf.RunMode)
	_=r.Run(config.Conf.AdminHost+":"+config.Conf.AdminPort)
}