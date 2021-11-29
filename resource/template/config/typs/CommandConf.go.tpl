package typs

type CommandConf struct {
	Enable    bool                 `json:"enable"`
	Cronds    []*CommandConf_Crond `json:"cronds"`
	Residents []*Command_Resident  `json:"residents"`
}

type CommandConf_Crond struct {
	Cron    string   `json:"cron"`
	Command string   `json:"command"`
	Options []string `json:"options"`
}

type Command_Resident struct {
	Wait    int      `json:"wait"`
	Command string   `json:"command"`
	Options []string `json:"options"`
}

func (i *CommandConf) ConfigName() string {
	return "CommandConf"
}
