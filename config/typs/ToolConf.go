package typs

type ToolConf struct {
	Version string
}

func (i *ToolConf) ConfigName() string {
	return "ToolConf"
}
