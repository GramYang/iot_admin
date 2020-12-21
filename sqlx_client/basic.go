package sqlx_client

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"iot_admin/config"
)

var db *sqlx.DB

func SetUp(){
	dataSource:=fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",config.Conf.MysqlUserName,
		config.Conf.MysqlPassword,config.Conf.MysqlAddr,config.Conf.MysqlPort,config.Conf.MysqlDatabase)
	db=sqlx.MustConnect("mysql",dataSource)
}