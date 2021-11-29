package vars

import "framework/config/typs"

var Command = &typs.CommandConf{
	Enable: true,
	Cronds: []*typs.CommandConf_Crond{
		//{Cron: "* * * * * *", Command: "test::app", Options: []string{"--name=cron"}},
	},
	Residents: []*typs.Command_Resident{
		{Wait: 0, Command: "test::app", Options: []string{"--name=resident"}},
		//&Command_Resident{Wait: 0, Command: "test::test", Options: []string{}},
	},
}
