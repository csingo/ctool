package command

import (
	"fmt"
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
	"regexp"
	"strings"
)

type AppCommand struct{}

type protoService struct {
	Name string
	Rpc  []protoServiceRpc
}

type protoServiceRpc struct {
	Name string
	Req  string
	Rsp  string
}

func (i *AppCommand) Help() *cCommand.CommandHelpDoc {
	return &cCommand.CommandHelpDoc{
		CommandDesc: cCommand.CommandDesc{Name: "app", Desc: "应用"},
		MethodDesc: []cCommand.MethodDesc{
			{Name: "create", Desc: "创建app", Options: []string{"name"}},
			{Name: "controller", Desc: "创建controller", Options: []string{"app", "name"}},
			{Name: "command", Desc: "创建command", Options: []string{"app", "name"}},
			{Name: "service", Desc: "创建service", Options: []string{"app", "protoPath"}},
		},
		OptionDesc: []cCommand.OptionDesc{
			{Name: "app", Desc: "应用"},
			{Name: "name", Desc: "名称"},
			{Name: "protoPath", Desc: "proto目录"},
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

func (i *AppCommand) Service(app cCommand.Option, protoPath cCommand.Option) {
	var services = []protoService{}

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

	// proto文件
	var protoc = []string{
		"protoc",
		fmt.Sprintf("--go_out=base/%s", app.Value),
		fmt.Sprintf("--go-grpc_out=base/%s", app.Value),
	}
	err = filepath.Walk(protoPath.Value, func(protoFile string, info fs.FileInfo, err error) error {
		fileExt := path.Ext(protoFile)
		if fileExt != ".proto" {
			return nil
		}

		protoc = append(protoc, protoFile)

		// 读取文件
		content, err := ioutil.ReadFile(protoFile)
		if err != nil {
			log.Fatal(err)
		}

		reg := regexp.MustCompile(`service *(\w+) *\{((\r\n|\n)* *rpc +(\w+)\( *(\w+) *\) *returns *\( *(\w+) *\)( *)((;)|(\{ *\})))+ *(\r\n|\n)* *\}`)
		matches := reg.FindAllStringSubmatch(string(content), -1)
		//log.Printf("%+V", matches)

		for _, match := range matches {
			matchService := match[0]
			//log.Println(matchService)
			serviceReg := regexp.MustCompile(`service +(\w+) *\{`)
			serviceMatches := serviceReg.FindAllStringSubmatch(matchService, -1)
			//log.Printf("%+V", serviceMatches[0][1])
			tempService := protoService{
				Name: serviceMatches[0][1],
				Rpc:  []protoServiceRpc{},
			}

			rpcReg := regexp.MustCompile(`rpc *(\w+) *\( *(\w+) *\) *returns *\( *(\w+) *\)`)
			rpcMatches := rpcReg.FindAllStringSubmatch(matchService, -1)
			//log.Printf("%+V", rpcMatches)
			for _, rpcMatch := range rpcMatches {
				tempService.Rpc = append(tempService.Rpc, protoServiceRpc{
					Name: rpcMatch[1],
					Req:  rpcMatch[2],
					Rsp:  rpcMatch[3],
				})
			}

			services = append(services, tempService)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("%+v", services)

	// 执行 protoc
	// protoc --go_out=base/app --go-grpc_out=base/app proto/app/hello_service.proto proto/app/enum.proto
	cmd := exec.Command(protoc[0], protoc[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}

	// 获取模板路径
	gopath := cHelper.EnvToString("GOPATH", "")
	tplPath := filepath.Clean(gopath + "/pkg/mod/gitee.com/csingo/ctool@" + vars.Tool.Version + "/resource/template")
	tplServiceFilePath := fmt.Sprintf("%s/base/app/service_http.pb.go.tpl", tplPath)
	tplRpcFilePath := fmt.Sprintf("%s/base/app/service_rpc.pb.go.tpl", tplPath)
	tplCallFilePath := fmt.Sprintf("%s/base/app/call.pb.go.tpl", tplPath)
	tplAppServiceFilePath := fmt.Sprintf("%s/app/service/Service.go.tpl", tplPath)
	tplAppServiceRpcFilePath := fmt.Sprintf("%s/app/service/Rpc.go.tpl", tplPath)

	serviceTplContent, err := ioutil.ReadFile(tplServiceFilePath)
	if err != nil {
		log.Fatal(err)
	}
	rpcTplContent, err := ioutil.ReadFile(tplRpcFilePath)
	if err != nil {
		log.Fatal(err)
	}
	appServiceTplContent, err := ioutil.ReadFile(tplAppServiceFilePath)
	if err != nil {
		log.Fatal(err)
	}
	appRpcTplContent, err := ioutil.ReadFile(tplAppServiceRpcFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// 生成 httprpc server
	for _, service := range services {
		servicePbFilePath := fmt.Sprintf("%s/base/%s/%s_http.pb.go", dir, app.Value, service.Name)

		content := string(serviceTplContent)
		content = strings.ReplaceAll(content, "##SERVICE##", service.Name)
		content = strings.ReplaceAll(content, "##APP##", app.Value)

		for _, rpc := range service.Rpc {
			subContent := string(rpcTplContent)
			subContent = strings.ReplaceAll(content, "##SERVICE##", service.Name)
			subContent = strings.ReplaceAll(content, "##RPC##", rpc.Name)
			subContent = strings.ReplaceAll(content, "##REQ##", rpc.Req)
			subContent = strings.ReplaceAll(content, "##RSP##", rpc.Rsp)

			content = content + subContent
		}

		// 写文件
		err = ioutil.WriteFile(servicePbFilePath, []byte(content), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	// 生成 call
	callContent, err := ioutil.ReadFile(tplCallFilePath)
	if err != nil {
		log.Fatal(err)
	}
	targetCallFilePath := filepath.Clean(fmt.Sprintf("%s/base/%s/call.pb.go", dir, app.Value))
	err = ioutil.WriteFile(targetCallFilePath, callContent, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// 创建 service 文件
	for _, service := range services {
		appServiceFilePath := fmt.Sprintf("%s/%s/service/%s.go", dir, app.Value, service.Name)

		appServiceContent := string(appServiceTplContent)
		appServiceContent = strings.ReplaceAll(appServiceContent, "##PROJECT##", project)
		appServiceContent = strings.ReplaceAll(appServiceContent, "##APP##", app.Value)
		appServiceContent = strings.ReplaceAll(appServiceContent, "##SERVICE##", service.Name)

		for _, rpc := range service.Rpc {
			appRpcContent := string(appRpcTplContent)
			appRpcContent = strings.ReplaceAll(appRpcContent, "##SERVICE##", service.Name)
			appRpcContent = strings.ReplaceAll(appRpcContent, "##RPC##", rpc.Name)
			appRpcContent = strings.ReplaceAll(appRpcContent, "##REQ##", rpc.Req)
			appRpcContent = strings.ReplaceAll(appRpcContent, "##Rsp##", rpc.Rsp)

			appServiceContent = appServiceContent + appRpcContent
		}

		// 写文件
		err = ioutil.WriteFile(appServiceFilePath, []byte(appServiceContent), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 创建 controller 文件
}

func (i *AppCommand) Controller(app cCommand.Option, name cCommand.Option) {
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
	tplFilePath := fmt.Sprintf("%s/app/controller/HomeController.go.tpl", tplPath)

	// 读取文件
	content, err := ioutil.ReadFile(tplFilePath)
	if err != nil {
		log.Fatal(err)
	}

	contentStr := string(content)
	contentStr = strings.ReplaceAll(contentStr, "##PROJECT##", project)
	contentStr = strings.ReplaceAll(contentStr, "##CONTROLLER##", name.Value)

	// 写文件
	targetFilePath := filepath.Clean(fmt.Sprintf("%s/%s/controller/%s.go", dir, app.Value, name.Value))
	content = []byte(contentStr)
	err = ioutil.WriteFile(targetFilePath, content, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// 更新autoload
	autoloadFilePath := filepath.Clean(fmt.Sprintf("%s/autoload/controller.go", dir))
	autoloadContent, err := ioutil.ReadFile(autoloadFilePath)
	if err != nil {
		log.Fatal(err)
	}

	importName := fmt.Sprintf("%scontroller", app.Value)
	autoloadContentStr := string(autoloadContent)
	autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "    //TODO:ImportController", fmt.Sprintf("    //TODO:InitController\n   %s \"%s/%s/controller\"", importName, project, app.Value))
	autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "    //TODO:InitController", fmt.Sprintf("    //TODO:InitController\n    cServer.Inject(&%s.%s{})", importName, name.Value))

	// 写文件
	autoloadContent = []byte(autoloadContentStr)
	err = ioutil.WriteFile(autoloadFilePath, autoloadContent, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *AppCommand) Command(app cCommand.Option, name cCommand.Option) {
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
	tplFilePath := fmt.Sprintf("%s/app/command/TestCommand.go.tpl", tplPath)

	// 读取文件
	content, err := ioutil.ReadFile(tplFilePath)
	if err != nil {
		log.Fatal(err)
	}

	contentStr := string(content)
	contentStr = strings.ReplaceAll(contentStr, "##PROJECT##", project)
	contentStr = strings.ReplaceAll(contentStr, "##COMMAND##", name.Value)
	contentStr = strings.ReplaceAll(contentStr, "##COMMANDNAME##", strings.ToLower(name.Value))

	// 写文件
	targetFilePath := filepath.Clean(fmt.Sprintf("%s/%s/controller/%s.go", dir, app.Value, name.Value))
	content = []byte(contentStr)
	err = ioutil.WriteFile(targetFilePath, content, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// 更新autoload
	autoloadFilePath := filepath.Clean(fmt.Sprintf("%s/autoload/command.go", dir))
	autoloadContent, err := ioutil.ReadFile(autoloadFilePath)
	if err != nil {
		log.Fatal(err)
	}

	importName := fmt.Sprintf("%command", app.Value)
	autoloadContentStr := string(autoloadContent)
	autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "    //TODO:ImportCommand", fmt.Sprintf("    //TODO:ImportCommand\n   %s \"%s/%s/command\"", importName, project, app.Value))
	autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "    //TODO:InitCommand", fmt.Sprintf("    //TODO:InitCommand\n    cServer.Inject(&%s.%s{})", importName, name.Value))

	// 写文件
	autoloadContent = []byte(autoloadContentStr)
	err = ioutil.WriteFile(autoloadFilePath, autoloadContent, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *AppCommand) Orm(app cCommand.Option, name cCommand.Option) {}

func (i *AppCommand) Sdk(sdk cCommand.Option, protoPath cCommand.Option) {}
