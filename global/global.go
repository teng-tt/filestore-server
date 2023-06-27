package global

import (
	"database/sql"
	"filestore-server/config"
)

var (
	CONF   config.Server
	DBConn *sql.DB
)
