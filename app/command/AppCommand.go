package command

import (
	"gitee.com/csingo/ctool/config/vars"
	"gitee.com/csingo/ctool/core/cCommand"
	"gitee.com/csingo/ctool/core/cHelper"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type AppCommand struct{}

func (i *AppCommand) Help() *cCommand.CommandHelpDoc {
	return &cCommand.CommandHelpDoc{
		CommandDesc: cCommand.CommandDesc{Name: "app", Desc: "应用"},
		MethodDesc: []cCommand.MethodDesc{
			{Name: "create", Desc: "创建", Options: []string{"name"}},
		},
		OptionDesc: []cCommand.OptionDesc{
			{Name: "name", Desc: "名称"},
		},
	}
}

func (i *AppCommand) Create(name cCommand.Option) {
	// 获取配置
	files := vars.Tool.WriteFiles["app::create"]
	length := len(files)
	for i := 0; i < length; i++ {
		files[i] = filepath.Clean(files[i])
	}

	// 获取创建项目的路径
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// 读取project name
	//modFile := dir + "/go.mod"
	//modContent, err := ioutil.ReadFile(modFile)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//modArr := strings.Split(string(modContent), "\n")
	//if len(modArr) <= 0 {
	//	log.Fatal("mod file is empty")
	//}
	//mod := strings.Trim(strings.Trim(modArr[0], "module"), " ")

	// 获取模板路径
	gopath := cHelper.EnvToString("GOPATH", "")
	tplPath := filepath.Clean(gopath + "/pkg/mod/gitee.com/csingo/ctool@" + vars.Tool.Version + "/resource/template")

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
				tempPath = strings.Replace(tempPath, "app", name.Value, 1)
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
}
