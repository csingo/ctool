package cConfig

import "log"

func initConfigCenter() {
	enable, err := GetConf("ConfigCenterConf.Enable")
	if err != nil {
		enable = false
	}
	ConfigContainer.ConfigCenter.Enable = enable.(bool)

	driver, err := GetConf("ConfigCenterConf.Driver")
	if err != nil {
		driver = "nacos"
	}
	ConfigContainer.ConfigCenter.Driver = driver.(string)

	interval, err := GetConf("ConfigCenterConf.Interval")
	if err != nil {
		interval = 1
	}
	ConfigContainer.UpdateInterval = interval.(int)

	if ConfigContainer.ConfigCenter.Enable {
		switch ConfigContainer.ConfigCenter.Driver {
		case "nacos":
			initNacos()
			listenConfigCenter()
		default:
			log.Printf("cConfig center driver is not used : %+v", ConfigContainer.ConfigCenter.Driver)
			ConfigContainer.ConfigCenter.Enable = false
		}
	}
}

func listenConfigCenter() {
	names := []string{}
	for name, _ := range ConfigContainer.ConfigCenter.Listeners {
		names = append(names, name)
	}

	ListenConf(names, nil)
}

func readConfigCenter(data, group string) ([]byte, error) {
	var res []byte
	var err error
	switch ConfigContainer.ConfigCenter.Driver {
	default:
		res = nil
		err = nil
	case "nacos":
		res, err = readNacos(data, group)
	}

	return res, err
}
