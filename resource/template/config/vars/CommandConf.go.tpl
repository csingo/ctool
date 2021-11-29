package vars

import "gitee.com/csingo/ctool/config/typs"

var Command = &typs.CommandConf{
	Enable: true,
	Cronds: []*typs.CommandConf_Crond{
	},
	Residents: []*typs.Command_Resident{
	},
}
