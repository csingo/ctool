package vars

import "gitee.com/csingo/ctool/config/typs"

var Command = &typs.CommandConf{
	Enable: true,
	Cronds: []*typs.CommandConf_Crond{
		//{Cron: "* * * * * *", Command: "test::hello", Options: []string{"--name=cron"}},
	},
	Residents: []*typs.Command_Resident{
		{Wait: 0, Command: "test::hello", Options: []string{"--name=resident"}},
		//&Command_Resident{Wait: 0, Command: "test::test", Options: []string{}},
	},
}
