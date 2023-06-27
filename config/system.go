package config

type System struct {
	Addr string `json:"addr" yaml:"addr"`
}

type DBConfig struct {
	Host       string `json:"host" yaml:"host"`
	Port       string `json:"port" yaml:"port"`
	Database   string `json:"database" yaml:"database"`
	UserName   string `json:"userName" yaml:"userName"`
	DBPassword string `json:"DBPassword" yaml:"DBPassword"`
}
