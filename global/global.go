package global

import (
	"database/sql"
	"filestore-server/config"
	"github.com/garyburd/redigo/redis"
)

var (
	CONF      config.Server
	DBConn    *sql.DB
	CacheConn *redis.Pool
)

const (
	PWD_SALT = "#990"
)
