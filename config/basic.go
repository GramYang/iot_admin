package config

import (
	"encoding/json"
	"flag"
	g "github.com/GramYang/gylog"
	"io/ioutil"
)

var Conf = &Config{}

type Config struct {
	AdminHost     string `json:"admin_host"`
	AdminPort     string `json:"admin_port"`
	RunMode       string `json:"run_mode"`
	MysqlAddr     string `json:"mysql_addr"`
	MysqlPort     string `json:"mysql_port"`
	MysqlUserName string `json:"mysql_username"`
	MysqlPassword string `json:"mysql_password"`
	MysqlDatabase string `json:"mysql_database"`
}

func (c *Config) defaultConfig() {
	c.AdminHost = "127.0.0.1"
	c.AdminPort = "8086"
	c.RunMode = "debug"
	c.MysqlAddr = "106.54.87.204"
	c.MysqlPort = "3306"
	c.MysqlUserName = "cqdq"
	c.MysqlPassword = "cqdq12345"
	c.MysqlDatabase = "iot_admin"
}

func (c *Config) initConfig(path string) {
	c.defaultConfig()
	if path != "" {
		file, err := ioutil.ReadFile(path)
		if err != nil {
			g.Errorln(err)
		}
		if err = json.Unmarshal(file, c); err != nil {
			g.Errorln(err)
		}
	}
}

func SetUp() {
	var p string
	flag.StringVar(&p, "c", "", "配置文件路径")
	flag.Parse()
	Conf.initConfig(p)
}
