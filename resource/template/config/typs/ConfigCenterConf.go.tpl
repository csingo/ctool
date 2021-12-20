package typs

type ConfigCenterConf struct {
	Enable   bool                       `json:"enable"`
	Driver   string                     `json:"driver"`
	Interval int                        `json:"interval"`
	Channels *ConfigCenterConf_Channels `json:"channels"`
}

type ConfigCenterConf_Channels struct {
	Nacos  *ConfigCenterConf_NacosChannel  `json:"nacos"`
	Apollo *ConfigCenterConf_ApolloChannel `json:"apollo"`
}

type ConfigCenterConf_NacosChannel struct {
	Client    *ConfigCenterConf_NacosChannel_Client     `json:"client"`
	Listeners []*ConfigCenterConf_NacosChannel_Listener `json:"listens"`
}

type ConfigCenterConf_ApolloChannel struct{}

type ConfigCenterConf_NacosChannel_Client struct {
	Scheme    string `json:"scheme"`
	Host      string `json:"host"`
	Port      uint64 `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Namespace string `json:"namespace"`
	Path      string `json:"path"`
}

type ConfigCenterConf_NacosChannel_Listener struct {
	Group string `json:"group"`
	Data  string `json:"data"`
	Conf  string `json:"conf"`
}

func (i *ConfigCenterConf) ConfigName() string {
	return "ConfigCenterConf"
}
