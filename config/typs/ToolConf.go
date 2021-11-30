package typs

type ToolConf struct {
	Version    string
	WriteFiles map[string][]string
}

func (i *ToolConf) ConfigName() string {
	return "ToolConf"
}
