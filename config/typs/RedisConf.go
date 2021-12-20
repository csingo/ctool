package typs

type RedisConf struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Auth string `json:"auth"`
	Db   int    `json:"db"`
}

func (i *RedisConf) ConfigName() string {
	return "RedisConf"
}
