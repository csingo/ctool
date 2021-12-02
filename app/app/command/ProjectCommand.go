package command

import (
	"gitee.com/csingo/ctool/config/vars"
	"gitee.com/csingo/ctool/core/cCommand"
	"gitee.com/csingo/ctool/core/cHelper"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
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

	// 获取模板路径
	modcache := cHelper.GetGOENV("GOMODCACHE")
	tplPath := filepath.Clean(modcache + "/gitee.com/csingo/ctool@" + vars.Tool.Version + "/resource/template")

	// 遍历文件夹复制文件
	err = filepath.Walk(tplPath, func(filePath string, info fs.FileInfo, err error) error {
		state, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if !state.IsDir() {
			fileExt := path.Ext(filePath)
			tempPath := strings.TrimPrefix(filePath, tplPath)
			if cHelper.InArrayString(filepath.Clean(tempPath), files) && fileExt == ".tpl" {
				tempFilePath := filepath.Clean(dir + strings.TrimRight(tempPath, ".tpl"))
				tempDir := filepath.Dir(tempFilePath)
				err = os.MkdirAll(tempDir, 0755)
				if err != nil {
					return err
				}
				if cHelper.IsExistsPath(tempFilePath) {
					return nil
				}
				content, err := ioutil.ReadFile(filePath)
				if err != nil {
					return err
				}

				// 写文件
				contentStr := string(content)
				contentStr = strings.ReplaceAll(contentStr, "##PROJECT##", name.Value)
				content = []byte(contentStr)
				err = ioutil.WriteFile(tempFilePath, content, 0755)
				if err != nil {
					return err
				}

				log.Println(tempDir)
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = shell([]string{
		"go mod init "+name.Value,
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
