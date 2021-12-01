package typs

type LogConf struct {
	Level int
}

func (i *LogConf) ConfigName() string {
	return "LogConf"
}
