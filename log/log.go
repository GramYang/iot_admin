package log

import (
	g "github.com/GramYang/gylog"
	"os"
)

//不使用gylog的日志分割
func InitLog2(mode int) {
	g.SetFlags(g.Lshortfile)
	g.SetLevel(g.LevelDebug)
	switch mode {
	case 0:
		file, err := os.OpenFile("iot_admin_log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		g.SetOutput(file)
	case 1:
		g.SetOutput(os.Stderr)
	}
}
