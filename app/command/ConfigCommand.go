package command

import (
	"fmt"
	"gitee.com/csingo/ctool/config/vars"
	"gitee.com/csingo/ctool/core/cCommand"
	"gitee.com/csingo/ctool/core/cHelper"
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

	// 读取project name
	modFile := dir + "/go.mod"
	modContent, err := ioutil.ReadFile(modFile)
	if err != nil {
		log.Fatal(err)
	}
	modArr := strings.Split(string(modContent), "\n")
	if len(modArr) <= 0 {
		log.Fatal("mod file is empty")
	}
	project := strings.Trim(strings.Trim(modArr[0], "module"), " ")

	// 获取模板路径
	gopath := cHelper.EnvToString("GOPATH", "")
	tplPath := filepath.Clean(gopath + "/pkg/mod/gitee.com/csingo/ctool@" + vars.Tool.Version + "/resource/template")
	tplAutoloadConfigFilePath := fmt.Sprintf("%s/autoload/config.go.tpl", tplPath)
	tplConfigTypFilePath := fmt.Sprintf("%s/config/typs/TestConf.go.tpl", tplPath)
	tplConfigVarFilePath := fmt.Sprintf("%s/config/vars/TestConf.go.tpl", tplPath)
	targetConfigTypFilePath := fmt.Sprintf("%s/config/typs/%s.go", dir, fullConfName)
	targetConfigVarFilePath := fmt.Sprintf("%s/config/vars/%s.go", dir, fullConfName)

	// 读取typs文件
	typsContent, err := ioutil.ReadFile(tplConfigTypFilePath)
	if err != nil {
		log.Fatal(err)
	}
	typsContentStr := string(typsContent)
	typsContentStr = strings.ReplaceAll(typsContentStr, "##CONFIG##", fullConfName)
	// 写typs文件
	if !cHelper.IsExistsPath(targetConfigTypFilePath) {
		err = ioutil.WriteFile(targetConfigTypFilePath, []byte(typsContentStr), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 读取vars文件
	varsContent, err := ioutil.ReadFile(tplConfigVarFilePath)
	if err != nil {
		log.Fatal(err)
	}
	varsContentStr := string(varsContent)
	varsContentStr = strings.ReplaceAll(varsContentStr, "##PROJECT##", project)
	varsContentStr = strings.ReplaceAll(varsContentStr, "##CONFIG##", fullConfName)
	varsContentStr = strings.ReplaceAll(varsContentStr, "##CONFIGNAME##", confName)
	// 写vars文件
	if !cHelper.IsExistsPath(targetConfigVarFilePath) {
		err = ioutil.WriteFile(targetConfigVarFilePath, []byte(varsContentStr), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 读取autoload文件
	autoloadContent, err := ioutil.ReadFile(tplAutoloadConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}
	autoloadContentStr := string(autoloadContent)
	if !strings.Contains(autoloadContentStr, confName) {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "//TODO:InitConfig", fmt.Sprintf("%s\tcServer.Inject(vars.%s)", "//TODO:InitConfig", confName))
		// 写autoload
		targetAutoloadFilePath := filepath.Clean(fmt.Sprintf("%s/autoload/config.go", dir))
		err = ioutil.WriteFile(targetAutoloadFilePath, []byte(autoloadContentStr), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}
