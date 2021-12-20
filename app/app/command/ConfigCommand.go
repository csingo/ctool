package command

import (
	"fmt"
	"gitee.com/csingo/ctool/core/cCommand"
	"gitee.com/csingo/ctool/core/cHelper"
	"gitee.com/csingo/ctool/resource/asset"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ConfigCommand struct{}

func (i *ConfigCommand) Help() *cCommand.CommandHelpDoc {
	return &cCommand.CommandHelpDoc{
		CommandDesc: cCommand.CommandDesc{Name: "config", Desc: "配置"},
		MethodDesc: []cCommand.MethodDesc{
			{Name: "create", Desc: "生成", Options: []string{"name"}},
		},
		OptionDesc: []cCommand.OptionDesc{
			{Name: "name", Desc: "名称"},
		},
	}
}

func (i *ConfigCommand) Create(name cCommand.Option) {
	fullConfName := name.Value
	confName := strings.TrimRight(fullConfName, "Conf")
	// 获取创建项目的路径
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// 读取go.mod
	project, err := cHelper.GetModName()
	if err != nil {
		log.Fatal(err)
	}

	// 读取模板
	autoloadContent, err := asset.Asset("resource/template/autoload/config.go.tpl")
	if err != nil {
		log.Fatal(err)
	}
	typsContent, err := asset.Asset("resource/template/config/typs/TestConf.go.tpl")
	if err != nil {
		log.Fatal(err)
	}
	varsContent, err := asset.Asset("resource/template/config/vars/TestConf.go.tpl")
	if err != nil {
		log.Fatal(err)
	}

	// 获取模板路径
	targetConfigTypFilePath := fmt.Sprintf("%s/config/typs/%s.go", dir, fullConfName)
	targetConfigVarFilePath := fmt.Sprintf("%s/config/vars/%s.go", dir, fullConfName)

	// 读取typs文件
	typsContentStr := cHelper.ReplaceAllFromMap(string(typsContent), map[string]string{
		"##CONFIG##": fullConfName,
	})
	// 写typs文件
	if !cHelper.IsExistsPath(targetConfigTypFilePath) {
		err = ioutil.WriteFile(targetConfigTypFilePath, []byte(typsContentStr), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 读取vars文件
	varsContentStr := cHelper.ReplaceAllFromMap(string(varsContent), map[string]string{
		"##PROJECT##":    project,
		"##CONFIG##":     fullConfName,
		"##CONFIGNAME##": confName,
	})
	// 写vars文件
	if !cHelper.IsExistsPath(targetConfigVarFilePath) {
		err = ioutil.WriteFile(targetConfigVarFilePath, []byte(varsContentStr), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 读取autoload文件
	autoloadContentStr := string(autoloadContent)
	if !strings.Contains(autoloadContentStr, confName) {
		autoloadContentStr = cHelper.ReplaceAllFromMap(autoloadContentStr, map[string]string{
			"##PROJECT##":       project,
			"//TODO:InitConfig": fmt.Sprintf("%s\n\tcServer.Inject(vars.%s)", "//TODO:InitConfig", confName),
		})
		// 写autoload
		targetAutoloadFilePath := filepath.Clean(fmt.Sprintf("%s/autoload/config.go", dir))
		err = ioutil.WriteFile(targetAutoloadFilePath, []byte(autoloadContentStr), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}
