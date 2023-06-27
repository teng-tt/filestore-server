package initiallize

import (
	"database/sql"
	"filestore-server/global"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func DBInit() {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		global.CONF.DBConfig.UserName,
		global.CONF.DBConfig.DBPassword,
		global.CONF.DBConfig.Host,
		global.CONF.DBConfig.Port,
		global.CONF.DBConfig.Database)
	
	db, _ := sql.Open("mysql", url)
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}
	global.DBConn = db
}
