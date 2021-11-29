package cConfig

import "log"

func initConfigCenter() {
	enable, err := GetConf("ConfigCenterConf.Enable")
	if err != nil {
		enable = false
	}
	configContainer.configCenter.enable = enable.(bool)

	driver, err := GetConf("ConfigCenterConf.Driver")
	if err != nil {
		driver = "nacos"
	}
	configContainer.configCenter.driver = driver.(string)

	interval, err := GetConf("ConfigCenterConf.Interval")
	if err != nil {
		interval = 1
	}
	configContainer.updateInterval = interval.(int)

	if configContainer.configCenter.enable {
		switch configContainer.configCenter.driver {
		case "nacos":
			initNacos()
			listenConfigCenter()
		default:
			log.Printf("cConfig center driver is not used : %+v", configContainer.configCenter.driver)
			configContainer.configCenter.enable = false
		}
	}
}

func listenConfigCenter() {
	names := []string{}
	for name, _ := range configContainer.configCenter.listeners {
		names = append(names, name)
	}

	ListenConf(names, nil)
}

func readConfigCenter(data, group string) ([]byte, error) {
	var res []byte
	var err error
	switch configContainer.configCenter.driver {
	default:
		res = nil
		err = nil
	case "nacos":
		res, err = readNacos(data, group)
	}

	return res, err
}
