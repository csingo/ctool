package typs

type ToolConf struct {
	Version               string
	ProtocDownload        *ToolConf_ProtocAddr
	ProtoGenGoPackage     string
	ProtoGenGoGrpcPackage string
	WriteFiles            map[string][]string
}

type ToolConf_ProtocAddr struct {
	Win   string
	Mac   string
	Linux string
}

func (i *ToolConf) ConfigName() string {
	return "ToolConf"
}
