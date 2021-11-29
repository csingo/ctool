package typs

type DatabaseConf struct {
	Driver   string                 `json:"driver"`
	Channels *DatabaseConf_Channels `json:"channels"`
}

type DatabaseConf_Channels struct {
	Default *DatabaseConf_MysqlChannel `json:"default"`
}

type DatabaseConf_MysqlChannel struct {
	Dsn             string `json:"dsn"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Database        string `json:"database"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Charset         string `json:"charset"`
	Parsetime       string `json:"parsetime"`
	Loc             string `json:"loc"`
	MaxIgleConns    int    `json:"max_igle_conns"`
	MaxOpenConns    int    `json:"max_open_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime"`
}

func (i *DatabaseConf) ConfigName() string {
	return "DatabaseConf"
}
