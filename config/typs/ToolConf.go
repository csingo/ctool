package typs

type ToolConf struct {
	CopyFiles    []string
	RewiteParams []string
}

func (i *ToolConf) ConfigName() string {
	return "ToolConf"
}
