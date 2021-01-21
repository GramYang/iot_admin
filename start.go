package main

import (
	"github.com/gin-gonic/gin"
	"iot_admin/config"
	"iot_admin/log"
	"iot_admin/router"
	sc "iot_admin/sqlx_client"
)

func main() {
	config.SetUp()
	if config.Conf.LocalDebug {
		log.InitLog2(1)
	} else { //log输出到文件，取消debug信息
		log.InitLog2(0)
		log.IsDebug(false)
	}
	sc.SetUp()
	r := router.NewRouter()
	gin.SetMode(config.Conf.RunMode)
	_ = r.Run(config.Conf.AdminHost + ":" + config.Conf.AdminPort)
}
