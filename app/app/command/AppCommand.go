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
			{Name: "middleware", Desc: "创建controller", Options: []string{"app", "name"}},
			{Name: "command", Desc: "创建command", Options: []string{"app", "name"}},
			{Name: "service", Desc: "创建service", Options: []string{"app", "protoPath"}},
		},
		OptionDesc: []cCommand.OptionDesc{
			{Name: "app", Desc: "应用"},
			{Name: "name", Desc: "名称"},
			{Name: "protoPath", Desc: "proto目录, 完整路径"},
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
				tempPath = strings.TrimPrefix(tempPath, "/app")
				tempPath = strings.Replace(tempPath, "app", name.Value, 1)
				tempPath = "/app" + tempPath
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

	// proto文件分析
	var outPath = filepath.Clean(fmt.Sprintf("base/%s", app.Value))
	err = os.MkdirAll(outPath, 0755)
	if err != nil {
		log.Fatal(err)
	}
	var protoc = []string{
		"protoc",
		fmt.Sprintf("--go_out=%s", outPath),
		fmt.Sprintf("--go-grpc_out=%s", outPath),
		"-I",
		filepath.Clean(protoPath.Value),
	}
	err = filepath.Walk(protoPath.Value, func(protoFile string, info fs.FileInfo, err error) error {
		fileExt := path.Ext(protoFile)
		filename := path.Base(protoFile)
		if fileExt != ".proto" {
			return nil
		}

		protoc = append(protoc, filename)

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

	// 执行 protoc
	// protoc --go_out=base/test --go-grpc_out=base/test -I D:\\Qdtech\\projects\\application-services\\ctool\\proto\\app enum.proto error.proto
	cmd := exec.Command(protoc[0], protoc[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}

	// 获取模板路径
	modcache := cHelper.GetGOENV("GOMODCACHE")
	tplPath := filepath.Clean(modcache + "/gitee.com/csingo/ctool@" + vars.Tool.Version + "/resource/template")
	tplServiceFilePath := fmt.Sprintf("%s/base/app/service_http.pb.go.tpl", tplPath)
	tplRpcFilePath := fmt.Sprintf("%s/base/app/service_rpc.pb.go.tpl", tplPath)
	tplCallFilePath := fmt.Sprintf("%s/base/app/call.pb.go.tpl", tplPath)
	tplAppServiceFilePath := fmt.Sprintf("%s/app/app/service/Service.go.tpl", tplPath)
	tplAppServiceRpcFilePath := fmt.Sprintf("%s/app/app/service/Rpc.go.tpl", tplPath)

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

	// 生成 httprpc server 和 http client
	for _, service := range services {
		servicePbFilePath := fmt.Sprintf("%s/base/%s/%s_%s_%s_http.pb.go", dir, app.Value, project, app.Value, service.Name)

		var contentByte = serviceTplContent
		var content string
		content = string(contentByte)
		content = strings.ReplaceAll(content, "##SERVICE##", service.Name)
		content = strings.ReplaceAll(content, "##APP##", app.Value)

		for _, rpc := range service.Rpc {
			var subContentByte = rpcTplContent
			var subContent string
			subContent = string(subContentByte)
			subContent = strings.ReplaceAll(subContent, "##APP##", app.Value)
			subContent = strings.ReplaceAll(subContent, "##SERVICE##", service.Name)
			subContent = strings.ReplaceAll(subContent, "##RPC##", rpc.Name)
			subContent = strings.ReplaceAll(subContent, "##REQ##", rpc.Req)
			subContent = strings.ReplaceAll(subContent, "##RSP##", rpc.Rsp)

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
	callContentStr := string(callContent)
	callContentStr = strings.ReplaceAll(callContentStr, "##PROJECT##", project)
	callContentStr = strings.ReplaceAll(callContentStr, "##APP##", app.Value)
	targetCallFilePath := filepath.Clean(fmt.Sprintf("%s/base/%s/call.pb.go", dir, app.Value))
	if !cHelper.IsExistsPath(targetCallFilePath) {
		err = ioutil.WriteFile(targetCallFilePath, []byte(callContentStr), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 创建 service 文件
	for _, service := range services {
		appServiceFilePath := fmt.Sprintf("%s/app/%s/service/%s.go", dir, app.Value, service.Name)
		if cHelper.IsExistsPath(appServiceFilePath) {
			appServiceTplContent, err = ioutil.ReadFile(appServiceFilePath)
			if err != nil {
				log.Fatalln(err)
			}
		}

		appServiceContent := string(appServiceTplContent)
		appServiceContent = strings.ReplaceAll(appServiceContent, "##PROJECT##", project)
		appServiceContent = strings.ReplaceAll(appServiceContent, "##APP##", app.Value)
		appServiceContent = strings.ReplaceAll(appServiceContent, "##SERVICE##", service.Name)

		for _, rpc := range service.Rpc {
			appRpcContent := string(appRpcTplContent)
			appRpcContent = strings.ReplaceAll(appRpcContent, "##SERVICE##", service.Name)
			appRpcContent = strings.ReplaceAll(appRpcContent, "##APP##", app.Value)
			appRpcContent = strings.ReplaceAll(appRpcContent, "##RPC##", rpc.Name)
			appRpcContent = strings.ReplaceAll(appRpcContent, "##REQ##", rpc.Req)
			appRpcContent = strings.ReplaceAll(appRpcContent, "##RSP##", rpc.Rsp)

			rpcReg := regexp.MustCompile(`func *\(\w+ *\*` + service.Name + `\) *` + rpc.Name + ` *\(`)
			if !rpcReg.MatchString(appServiceContent) {
				appServiceContent = appServiceContent + appRpcContent
			}
		}

		// 写文件
		err = ioutil.WriteFile(appServiceFilePath, []byte(appServiceContent), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	// 执行 go mod tidy
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}

	// 更新 autoload
	autoloadFilePath := filepath.Clean(fmt.Sprintf("%s/autoload/service.go", dir))
	autoloadContent, err := ioutil.ReadFile(autoloadFilePath)
	if err != nil {
		log.Fatal(err)
	}

	importName := fmt.Sprintf("%sService", app.Value)
	autoloadContentStr := string(autoloadContent)
	if !strings.Contains(autoloadContentStr, "import (") {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "func initService() {", "import (\n)\n\nfunc initService() {")
	}
	if !strings.Contains(autoloadContentStr, "//TODO:ImportService") {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "import (", "import (\n    //TODO:ImportService")
	}
	if !strings.Contains(autoloadContentStr, "//TODO:InitService") {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "func initService() {", "func initService() {\n    //TODO:InitService")
	}
	if !strings.Contains(autoloadContentStr, "core/cServer") {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "    //TODO:ImportService", fmt.Sprintf("    //TODO:ImportService\n   \"%s/core/cServer\"", project))
	}
	if !strings.Contains(autoloadContentStr, importName) {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "    //TODO:ImportService", fmt.Sprintf("    //TODO:ImportService\n   %s \"%s/app/%s/service\"", importName, project, app.Value))
	}
	for _, service := range services {
		var serviceName = fmt.Sprintf("%s.%s", importName, service.Name)
		if !strings.Contains(autoloadContentStr, serviceName) {
			autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "    //TODO:InitService", fmt.Sprintf("    //TODO:InitService\n    cServer.Inject(&%s{})", serviceName))
		}
	}

	// 写文件
	autoloadContent = []byte(autoloadContentStr)
	err = ioutil.WriteFile(autoloadFilePath, autoloadContent, 0755)
	if err != nil {
		log.Fatal(err)
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
	modcache := cHelper.GetGOENV("GOMODCACHE")
	tplPath := filepath.Clean(modcache + "/gitee.com/csingo/ctool@" + vars.Tool.Version + "/resource/template")
	tplFilePath := fmt.Sprintf("%s/app/app/controller/HomeController.go.tpl", tplPath)

	// 读取文件
	content, err := ioutil.ReadFile(tplFilePath)
	if err != nil {
		log.Fatal(err)
	}

	contentStr := string(content)
	contentStr = strings.ReplaceAll(contentStr, "##PROJECT##", project)
	contentStr = strings.ReplaceAll(contentStr, "##CONTROLLER##", name.Value)

	// 写文件
	targetFilePath := filepath.Clean(fmt.Sprintf("%s/app/%s/controller/%s.go", dir, app.Value, name.Value))
	if cHelper.IsExistsPath(targetFilePath) {
		log.Fatalf("file is exists: %s", targetFilePath)
	}
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
	if !strings.Contains(autoloadContentStr, "import (") {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "func initController() {", "import (\n)\n\nfunc initController() {")
	}
	if !strings.Contains(autoloadContentStr, "//TODO:ImportController") {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "import (", "import (\n    //TODO:ImportController")
	}
	if !strings.Contains(autoloadContentStr, "//TODO:InitController") {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "func initController() {", "func initController() {\n    //TODO:InitController")
	}
	if !strings.Contains(autoloadContentStr, "core/cServer") {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "    //TODO:ImportController", fmt.Sprintf("    //TODO:ImportController\n   \"%s/core/cServer\"", project))
	}
	if !strings.Contains(autoloadContentStr, importName) {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "    //TODO:ImportController", fmt.Sprintf("    //TODO:ImportController\n   %s \"%s/%s/controller\"", importName, project, app.Value))
	}
	controllerName := fmt.Sprintf("%s.%s", importName, name.Value)
	if !strings.Contains(autoloadContentStr, controllerName) {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "    //TODO:InitController", fmt.Sprintf("    //TODO:InitController\n    cServer.Inject(&%s{})", controllerName))
	}

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
	modcache := cHelper.GetGOENV("GOMODCACHE")
	tplPath := filepath.Clean(modcache + "/gitee.com/csingo/ctool@" + vars.Tool.Version + "/resource/template")
	tplFilePath := fmt.Sprintf("%s/app/app/command/TestCommand.go.tpl", tplPath)

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
	targetFilePath := filepath.Clean(fmt.Sprintf("%s/app/%s/command/%s.go", dir, app.Value, name.Value))
	if cHelper.IsExistsPath(targetFilePath) {
		log.Fatalf("file is exists: %s", targetFilePath)
	}
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

	importName := fmt.Sprintf("%scommand", app.Value)
	autoloadContentStr := string(autoloadContent)
	if !strings.Contains(autoloadContentStr, "import (") {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "func initCommand() {", "import (\n)\n\nfunc initCommand() {")
	}
	if !strings.Contains(autoloadContentStr, "//TODO:ImportCommand") {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "import (", "import (\n    //TODO:ImportCommand")
	}
	if !strings.Contains(autoloadContentStr, "//TODO:InitCommand") {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "func initCommand() {", "func initCommand() {\n    //TODO:InitCommand")
	}
	if !strings.Contains(autoloadContentStr, "core/cServer") {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "    //TODO:ImportCommand", fmt.Sprintf("    //TODO:ImportCommand\n   \"%s/core/cServer\"", project))
	}
	if !strings.Contains(autoloadContentStr, importName) {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "    //TODO:ImportCommand", fmt.Sprintf("    //TODO:ImportCommand\n   %s \"%s/%s/command\"", importName, project, app.Value))
	}
	commandName := fmt.Sprintf("%s.%s", importName, name.Value)
	if !strings.Contains(autoloadContentStr, commandName) {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "    //TODO:InitCommand", fmt.Sprintf("    //TODO:InitCommand\n    cServer.Inject(&%s{})", commandName))
	}

	// 写文件
	autoloadContent = []byte(autoloadContentStr)
	err = ioutil.WriteFile(autoloadFilePath, autoloadContent, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *AppCommand) Middleware(app cCommand.Option, name cCommand.Option) {
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
	modcache := cHelper.GetGOENV("GOMODCACHE")
	tplPath := filepath.Clean(modcache + "/gitee.com/csingo/ctool@" + vars.Tool.Version + "/resource/template")
	tplFilePath := fmt.Sprintf("%s/app/app/middleware/TestMiddleware.go.tpl", tplPath)

	// 读取文件
	content, err := ioutil.ReadFile(tplFilePath)
	if err != nil {
		log.Fatal(err)
	}

	contentStr := string(content)
	contentStr = strings.ReplaceAll(contentStr, "##PROJECT##", project)
	contentStr = strings.ReplaceAll(contentStr, "##MIDDLEWARE##", name.Value)

	// 写文件
	targetFilePath := filepath.Clean(fmt.Sprintf("%s/app/%s/middleware/%s.go", dir, app.Value, name.Value))
	if cHelper.IsExistsPath(targetFilePath) {
		log.Fatalf("file is exists: %s", targetFilePath)
	}
	content = []byte(contentStr)
	err = ioutil.WriteFile(targetFilePath, content, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// 更新autoload
	autoloadFilePath := filepath.Clean(fmt.Sprintf("%s/autoload/middleware.go", dir))
	autoloadContent, err := ioutil.ReadFile(autoloadFilePath)
	if err != nil {
		log.Fatal(err)
	}

	importName := fmt.Sprintf("%smiddleware", app.Value)
	autoloadContentStr := string(autoloadContent)
	if !strings.Contains(autoloadContentStr, "import (") {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "func initMiddleware() {", "import (\n)\n\nfunc initMiddleware() {")
	}
	if !strings.Contains(autoloadContentStr, "//TODO:ImportMiddleware") {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "import (", "import (\n    //TODO:ImportMiddleware")
	}
	if !strings.Contains(autoloadContentStr, "//TODO:InitMiddleware") {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "func initMiddleware() {", "func initMiddleware() {\n    //TODO:InitMiddleware")
	}
	if !strings.Contains(autoloadContentStr, "core/cServer") {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "    //TODO:ImportMiddleware", fmt.Sprintf("    //TODO:ImportMiddleware\n   \"%s/core/cServer\"", project))
	}
	if !strings.Contains(autoloadContentStr, importName) {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "    //TODO:ImportMiddleware", fmt.Sprintf("    //TODO:ImportMiddleware\n   %s \"%s/%s/middleware\"", importName, project, app.Value))
	}
	middlewareName := fmt.Sprintf("%s.%s", importName, name.Value)
	if !strings.Contains(autoloadContentStr, middlewareName) {
		autoloadContentStr = strings.ReplaceAll(autoloadContentStr, "    //TODO:InitMiddleware", fmt.Sprintf("    //TODO:InitMiddleware\n    cServer.Inject(&%s{})", middlewareName))
	}

	// 写文件
	autoloadContent = []byte(autoloadContentStr)
	err = ioutil.WriteFile(autoloadFilePath, autoloadContent, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *AppCommand) Orm(app cCommand.Option, name cCommand.Option) {}

