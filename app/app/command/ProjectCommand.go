package command

import (
	"fmt"
	"gitee.com/csingo/ctool/config/vars"
	"gitee.com/csingo/ctool/core/cCommand"
	"gitee.com/csingo/ctool/core/cHelper"
	"gitee.com/csingo/ctool/resource/asset"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ProjectCommand struct{}

func (i *ProjectCommand) Help() *cCommand.CommandHelpDoc {
	return &cCommand.CommandHelpDoc{
		CommandDesc: cCommand.CommandDesc{Name: "project", Desc: "项目"},
		MethodDesc: []cCommand.MethodDesc{
			{Name: "create", Desc: "创建", Options: []string{"name"}},
		},
		OptionDesc: []cCommand.OptionDesc{
			{Name: "name", Desc: "名称"},
		},
	}
}

func (i *ProjectCommand) Create(name cCommand.Option) {
	// 获取配置
	files := vars.Tool.WriteFiles["project::create"]
	length := len(files)
	for i := 0; i < length; i++ {
		files[i] = filepath.Clean(files[i])
	}

	// 获取创建项目的路径
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		content, err := asset.Asset(f)
		if err != nil {
			continue
		}

		tempFilePath := strings.TrimPrefix(f, "resource/template")
		tempFilePath = strings.TrimRight(tempFilePath, ".tpl")
		tempFilePath = fmt.Sprintf("%s%s", dir, tempFilePath)

		tempDir := filepath.Dir(tempFilePath)
		err = os.MkdirAll(tempDir, 0755)
		if err != nil {
			continue
		}
		if cHelper.IsExistsPath(tempFilePath) {
			continue
		}

		// 写文件
		contentStr := string(content)
		contentStr = cHelper.ReplaceAllFromMap(contentStr, map[string]string{"##PROJECT##": name.Value})
		content = []byte(contentStr)
		ioutil.WriteFile(tempFilePath, content, 0755)

		log.Println("write success:", tempFilePath)
	}

	err = shell([]string{
		"go mod init " + name.Value,
		"go mod tidy",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func shell(commands []string) error {

	for _, command := range commands {
		log.Println(command)
		commandArr := strings.Split(command, " ")
		cmd := exec.Command(commandArr[0], commandArr[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
