package main

import (
	"filestore-server/global"
	"filestore-server/initiallize"
)

func main() {
	initiallize.Viper()
	r := initiallize.Routers()
	initiallize.DBInit()
	panic(r.Run(global.CONF.System.Addr))
}
