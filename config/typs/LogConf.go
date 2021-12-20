package typs

type LogConf struct {
	Level  int    `json:"level"`
	Output string `json:"output"`
}

func (i *LogConf) ConfigName() string {
	return "LogConf"
}
