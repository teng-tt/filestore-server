package main

import (
	"filestore-server/global"
	"filestore-server/initiallize"
)

// 项目启动入口
func main() {
	r := initiallize.Routers()
	initiallize.Viper()
	initiallize.DBInit()
	initiallize.CacheInit()
	panic(r.Run(global.CONF.System.Addr))
}
