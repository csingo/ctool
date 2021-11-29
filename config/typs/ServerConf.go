package typs

type ServerConf struct {
	GoMod      string                 `json:"go_mod"`
	HttpServer *ServerConf_HttpServer `json:"http_server"`
}

type ServerConf_HttpServer struct {
	Enable         bool `json:"enable"`
	Port           int  `json:"port"`
	ReadTimeout    int  `json:"read_timeout"`
	WriteTimeout   int  `json:"write_timeout"`
	ExitTimeout    int  `json:"exit_timeout"`
	MaxHeaderBytes int  `json:"max_header_bytes"`
}

func (i *ServerConf) ConfigName() string {
	return "ServerConf"
}
