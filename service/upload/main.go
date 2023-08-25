package main

import (
	"filestore-server/global"
	"filestore-server/initiallize"
)

func main() {
	initiallize.Viper()
	initiallize.DBInit()
	initiallize.CacheInit()
	r := initiallize.Routers()
	panic(r.Run(global.CONF.System.Addr))
}
