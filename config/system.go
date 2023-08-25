package config

type System struct {
	Addr string `json:"addr" yaml:"addr"`
}

type MysqlConf struct {
	Host       string `json:"host" yaml:"host"`
	Port       string `json:"port" yaml:"port"`
	Database   string `json:"database" yaml:"database"`
	UserName   string `json:"userName" yaml:"userName"`
	DBPassword string `json:"DBPassword" yaml:"DBPassword"`
}

type RedisConf struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Password string `json:"password" yaml:"password"`
}

type OssConf struct {
	Bucket          string `json:"bucket" yaml:"bucket"`
	Endpoint        string `json:"endpoint" yaml:"endpoint"`
	AccessKeyID     string `json:"accessKeyID" yaml:"accessKeyID"`
	AccessKeySecret string `json:"accessKeySecret" yaml:"accessKeySecret"`
}

type RabbitConf struct {
	AsyncTransferEnable  bool   `json:"asyncTransferEnable" yaml:"asyncTransferEnable"`   // 是否开启文件异步转移默认同步
	RabbitURL            string `json:"rabbitURL" yaml:"rabbitURL"`                       // rabbitmq服务的入口
	TransExchangeName    string `json:"transExchangeName" yaml:"transExchangeName"`       // 用于文件transfer的交换机
	TransOSSQueueName    string `json:"transOSSQueueName" yaml:"transOSSQueueName"`       // oss转移队列名称
	TransOSSErrQueueName string `json:"transOSSErrQueueName" yaml:"transOSSErrQueueName"` // oss转移失败后写入另一个队列的队列名称
	TransOSSRoutingKey   string `json:"transOSSRoutingKey" yaml:"transOSSRoutingKey"`     // routingKey
}
