package config

type Server struct {
	System     System     `json:"system"`
	DBConf     MysqlConf  `json:"dbConf"`
	CacheConf  RedisConf  `json:"cacheConf"`
	OssConf    OssConf    `json:"ossConf"`
	RabbitConf RabbitConf `json:"rabbitConf"`
}
