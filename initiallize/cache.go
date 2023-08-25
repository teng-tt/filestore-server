package initiallize

import (
	"filestore-server/global"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

// 返回redis连接池
func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,                // 链接池最多有多少个可用的链接
		MaxActive:   30,                // 同时能够使用的链接
		IdleTimeout: 300 * time.Second, // 链接多久没使用就回收
		Dial: func() (redis.Conn, error) {
			// 打开链接
			c, err := redis.Dial("tcp", global.CONF.CacheConf.Host+global.CONF.CacheConf.Port)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			// 访问认证
			if _, err = c.Do("AUTH", global.CONF.CacheConf.Password); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		// 定时检查链接是否可用，不可用关闭链接
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}

func CacheInit() {
	pool := newRedisPool()
	global.CacheConn = pool
}
