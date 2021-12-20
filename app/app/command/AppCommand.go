package command

import (
	"fmt"
	"github.com/csingo/ctool/config/vars"
	"github.com/csingo/ctool/core/cCommand"
	"github.com/csingo/ctool/core/cHelper"
	"github.com/csingo/ctool/resource/asset"
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
	Name     string
	Rpc      []protoServiceRpc
	Packages []string
}

type protoEnum struct {
	Name   string
	Values map[int32]string
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

	// 读取go.mod
	project, err := cHelper.GetModName()
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		content, err := asset.Asset(f)
		if err != nil {
			continue
		}

		tempFilePath := strings.TrimPrefix(f, "resource/template/app/app")
		tempFilePath = strings.TrimRight(tempFilePath, ".tpl")
		tempFilePath = fmt.Sprintf("%s/app/%s/%s", dir, name.Value, tempFilePath)

		tempDir := filepath.Dir(tempFilePath)
		err = os.MkdirAll(tempDir, 0755)
		if err != nil {
			continue
		}
		if cHelper.IsExistsPath(tempFilePath) {
			continue
		}

		// 写文件
		contentStr := cHelper.ReplaceAllFromMap(string(content), map[string]string{
			"##PROJECT##": project,
			"##APP##":     name.Value,
		})
		err = ioutil.WriteFile(tempFilePath, []byte(contentStr), 0755)
		if err != nil {
			continue
		}

		log.Println("create success:", tempFilePath)
	}
}

func (i *AppCommand) Service(app cCommand.Option, protoPath cCommand.Option) {
	var services = []protoService{}
	var enums = []protoEnum{}

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

	// proto文件分析
	var outPath = filepath.Clean("base/cache")
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
	walkPath := filepath.Clean(fmt.Sprintf("%s/%s", protoPath.Value, app.Value))
	err = filepath.Walk(walkPath, func(protoFile string, info fs.FileInfo, err error) error {
		//log.Println(protoFile)
		fileExt := path.Ext(protoFile)
		filename := path.Base(protoFile)
		if fileExt != ".proto" {
			return nil
		}

		//protoc = append(protoc, filename)
		protoFileName := fmt.Sprintf("%s/%s", app.Value, filename)
		protoc = append(protoc, protoFileName)

		// 读取文件
		content, err := ioutil.ReadFile(protoFile)
		if err != nil {
			log.Fatal(err)
		}

		// 分析import文件
		importFiles := protoImportFiles(protoPath.Value, protoFileName)
		protoc = append(protoc, importFiles...)
		log.Println("proto import files: ", importFiles)

		// 获取service
		reg := regexp.MustCompile(`service *(\w+) *\{((\r\n|\n)* *rpc +(\w+) *\( *([\w\.]+) *\) *returns *\( *([\w\.]+) *\)( *)((;)|(\{ *\})))+ *(\r\n|\n)* *\}`)
		matches := reg.FindAllStringSubmatch(string(content), -1)
		//log.Printf("%+V", matches)

		for _, match := range matches {
			matchService := match[0]
			//log.Println(matchService)
			serviceReg := regexp.MustCompile(`service +(\w+) *\{`)
			serviceMatches := serviceReg.FindAllStringSubmatch(matchService, -1)
			//log.Printf("%+V", serviceMatches[0][1])
			tempService := protoService{
				Name:     serviceMatches[0][1],
				Rpc:      []protoServiceRpc{},
				Packages: protoImportPackage(protoPath.Value, protoFileName, project),
			}

			rpcReg := regexp.MustCompile(`rpc *(\w+) *\( *([\w\.]+) *\) *returns *\( *([\w\.]+) *\)`)
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

		// 获取enum
		enumReg := regexp.MustCompile(`enum +([\w_]+) *\{ *((\r\n|\n)* *[\w_]+ *= *([0-9]+); *(\/\/([^\r\n])*)*)* *(\r\n|\n)* *\}`)
		enumMatches := enumReg.FindAllStringSubmatch(string(content), -1)
		for _, e := range enumMatches {
			var tempProtoEnum = protoEnum{
				Name:   e[1],
				Values: map[int32]string{},
			}
			enumValueReg := regexp.MustCompile(` *([\w_]+) *= *([0-9]+) *; *(\/\/[^\r\n]*)*`)
			enumValueMatches := enumValueReg.FindAllStringSubmatch(e[0], -1)
			for _, v := range enumValueMatches {
				var enumValueDesc = v[3]
				if enumValueDesc == "" {
					enumValueDesc = v[1]
				}
				enumValueIndex := cHelper.ToInt32(v[2])
				tempProtoEnum.Values[enumValueIndex] = strings.Trim(strings.Trim(enumValueDesc, "/"), " ")
			}
			enums = append(enums, tempProtoEnum)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("load services: ", services)

	// 执行 protoc
	// protoc --go_out=base/test --go-grpc_out=base/test -I D:\\Qdtech\\projects\\application-services\\ctool\\proto\\app enum.proto error.proto
	cmd := exec.Command(protoc[0], protoc[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Println(protoc)
		log.Fatalln(err)
	}

	// 修改protoc生成的文件
	basePath := fmt.Sprintf("%s/base", dir)
	filepath.Walk(basePath, func(pbFile string, info fs.FileInfo, err error) error {
		if !strings.Contains(pbFile, ".pb.go") {
			return nil
		}

		data, _ := ioutil.ReadFile(pbFile)
		content := strings.ReplaceAll(string(data), "../", fmt.Sprintf("%s/base/", project))
		ioutil.WriteFile(pbFile, []byte(content), 0755)

		return nil
	})

	// 获取模板文件
	enumDescTplContent, err := asset.Asset("resource/template/base/app/enum_desc.go.tpl")
	if err != nil {
		log.Fatal(err)
	}
	serviceTplContent, err := asset.Asset("resource/template/base/app/service_http.pb.go.tpl")
	if err != nil {
		log.Fatal(err)
	}
	rpcTplContent, err := asset.Asset("resource/template/base/app/service_rpc.pb.go.tpl")
	if err != nil {
		log.Fatal(err)
	}
	appServiceTplContent, err := asset.Asset("resource/template/app/app/service/Service.go.tpl")
	if err != nil {
		log.Fatal(err)
	}
	appRpcTplContent, err := asset.Asset("resource/template/app/app/service/Rpc.go.tpl")
	if err != nil {
		log.Fatal(err)
	}
	callContent, err := asset.Asset("resource/template/base/app/call.pb.go.tpl")
	if err != nil {
		log.Fatal(err)
	}

	// 生成 enum_desc
	enumFilePath := fmt.Sprintf("%s/base/%s/enum_desc.go", dir, app.Value)
	enumContent := fmt.Sprintf("package %s\n", app.Value)
	for _, enumItem := range enums {
		var enumValueContent = "\n"
		for ei, ev := range enumItem.Values {
			enumValueContent = fmt.Sprintf("%s\t%d: \"%s\",\n", enumValueContent, ei, ev)
		}
		tempEnumContent := cHelper.ReplaceAllFromMap(string(enumDescTplContent), map[string]string{
			"##ENUM##":      enumItem.Name,
			"##ENUMVALUE##": enumValueContent,
		})
		enumContent = enumContent + tempEnumContent
	}
	// 写文件
	err = ioutil.WriteFile(enumFilePath, []byte(enumContent), 0755)
	if err != nil {
		log.Fatal(err)
	}

	// 生成 httprpc server 和 http client
	for _, service := range services {
		servicePbFilePath := fmt.Sprintf("%s/base/%s/%s_%s_%s_http.go", dir, app.Value, project, app.Value, service.Name)
		log.Println(servicePbFilePath)

		content := cHelper.ReplaceAllFromMap(string(serviceTplContent), map[string]string{
			"##SERVICE##": service.Name,
			"##APP##":     app.Value,
		})

		for _, importPackage := range service.Packages {
			content = cHelper.ReplaceAllFromMap(content, map[string]string{
				"//TODO:Import": fmt.Sprintf("%s\n\t\"%s\"", "//TODO:Import", importPackage),
			})
		}

		for _, rpc := range service.Rpc {
			subContent := cHelper.ReplaceAllFromMap(string(rpcTplContent), map[string]string{
				"##APP##":     app.Value,
				"##SERVICE##": service.Name,
				"##RPC##":     rpc.Name,
				"##REQ##":     rpc.Req,
				"##RSP##":     rpc.Rsp,
			})

			content = content + subContent
		}

		// 写文件
		err = ioutil.WriteFile(servicePbFilePath, []byte(content), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 生成 call
	callContentStr := cHelper.ReplaceAllFromMap(string(callContent), map[string]string{
		"##PROJECT##": project,
		"##APP##":     app.Value,
	})
	targetCallFilePath := filepath.Clean(fmt.Sprintf("%s/base/%s/call.go", dir, app.Value))
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

		appServiceContent := cHelper.ReplaceAllFromMap(string(appServiceTplContent), map[string]string{
			"##PROJECT##": project,
			"##APP##":     app.Value,
			"##SERVICE##": service.Name,
		})

		for _, importPackage := range service.Packages {
			appServiceContent = cHelper.ReplaceAllFromMap(appServiceContent, map[string]string{
				"//TODO:Import": fmt.Sprintf("%s\n\t\"%s\"", "//TODO:Import", importPackage),
			})
		}

		for _, rpc := range service.Rpc {
			reqApp := app.Value
			rspApp := app.Value
			req := rpc.Req
			rsp := rpc.Rsp

			if strings.Contains(rpc.Req, ".") {
				reqSplit := strings.Split(rpc.Req, ".")
				reqApp = reqSplit[0]
				req = reqSplit[1]
			}
			if strings.Contains(rpc.Rsp, ".") {
				rspSplit := strings.Split(rpc.Rsp, ".")
				rspApp = rspSplit[0]
				rsp = rspSplit[1]
			}

			appRpcContent := cHelper.ReplaceAllFromMap(string(appRpcTplContent), map[string]string{
				"##APP##":     app.Value,
				"##SERVICE##": service.Name,
				"##REQAPP##":  reqApp,
				"##RSPAPP##":  rspApp,
				"##RPC##":     rpc.Name,
				"##REQ##":     req,
				"##RSP##":     rsp,
			})

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
		log.Println(appServiceFilePath)
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
	autoloadContentStr := cHelper.ReplaceAllIfNotContain(string(autoloadContent), [][]string{
		{"import (", "func initService() {", "import (\n)\n\nfunc initService() {"},
		{"//TODO:ImportService", "import (", "import (\n\t//TODO:ImportService"},
		{"core/cServer", "//TODO:ImportService", fmt.Sprintf("//TODO:ImportService\n\t\"%s/core/cServer\"", project)},
		{importName, "//TODO:ImportService", fmt.Sprintf("//TODO:ImportService\n\t%s \"%s/app/%s/service\"", importName, project, app.Value)},
		{"//TODO:InitService", "func initService() {", "func initService() {\n\t//TODO:InitService"},
	})
	for _, service := range services {
		var serviceName = fmt.Sprintf("%s.%s", importName, service.Name)
		if !strings.Contains(autoloadContentStr, serviceName) {
			autoloadContentStr = cHelper.ReplaceAllIfNotContain(autoloadContentStr, [][]string{
				{serviceName, "//TODO:InitService", fmt.Sprintf("//TODO:InitService\n\tcServer.Inject(&%s{})", serviceName)},
			})
		}
		log.Println(serviceName)
	}

	// 写文件
	err = ioutil.WriteFile(autoloadFilePath, []byte(autoloadContentStr), 0755)
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

	// 读取go.mod
	project, err := cHelper.GetModName()
	if err != nil {
		log.Fatal(err)
	}

	// 读取文件
	content, err := asset.Asset("resource/template/app/app/controller/HomeController.go.tpl")
	if err != nil {
		log.Fatal(err)
	}

	contentStr := cHelper.ReplaceAllFromMap(string(content), map[string]string{
		"##PROJECT##":    project,
		"##CONTROLLER##": name.Value,
	})

	// 写文件
	targetFilePath := filepath.Clean(fmt.Sprintf("%s/app/%s/controller/%s.go", dir, app.Value, name.Value))
	if cHelper.IsExistsPath(targetFilePath) {
		log.Fatalf("file is exists: %s", targetFilePath)
	}
	err = ioutil.WriteFile(targetFilePath, []byte(contentStr), 0755)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("create success:", targetFilePath)

	// 更新autoload
	autoloadFilePath := filepath.Clean(fmt.Sprintf("%s/autoload/controller.go", dir))
	autoloadContent, err := ioutil.ReadFile(autoloadFilePath)
	if err != nil {
		log.Fatal(err)
	}

	importName := fmt.Sprintf("%scontroller", app.Value)
	controllerName := fmt.Sprintf("%s.%s", importName, name.Value)
	autoloadContentStr := cHelper.ReplaceAllIfNotContain(string(autoloadContent), [][]string{
		{"import (", "func initController() {", "import (\n)\n\nfunc initController() {"},
		{"//TODO:ImportController", "import (", "import (\n\t//TODO:ImportController"},
		{"//TODO:InitController", "func initController() {", "func initController() {\n\t//TODO:InitController"},
		{"core/cServer", "//TODO:ImportController", fmt.Sprintf("//TODO:ImportController\n\t\"%s/core/cServer\"", project)},
		{importName, "//TODO:ImportController", fmt.Sprintf("//TODO:ImportController\n\t%s \"%s/app/%s/controller\"", importName, project, app.Value)},
		{controllerName, "//TODO:InitController", fmt.Sprintf("//TODO:InitController\n\tcServer.Inject(&%s{})", controllerName)},
	})

	// 写文件
	err = ioutil.WriteFile(autoloadFilePath, []byte(autoloadContentStr), 0755)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("autoload success:", autoloadFilePath)
}

func (i *AppCommand) Command(app cCommand.Option, name cCommand.Option) {
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

	// 读取模板文件
	content, err := asset.Asset("resource/template/app/app/command/TestCommand.go.tpl")
	if err != nil {
		log.Fatal(err)
	}

	contentStr := cHelper.ReplaceAllFromMap(string(content), map[string]string{
		"##PROJECT##":     project,
		"##COMMAND##":     name.Value,
		"##COMMANDNAME##": strings.ToLower(name.Value),
	})

	// 写文件
	targetFilePath := filepath.Clean(fmt.Sprintf("%s/app/%s/command/%s.go", dir, app.Value, name.Value))
	if cHelper.IsExistsPath(targetFilePath) {
		log.Fatalf("file is exists: %s", targetFilePath)
	}
	err = ioutil.WriteFile(targetFilePath, []byte(contentStr), 0755)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("create success:", targetFilePath)

	// 更新autoload
	autoloadFilePath := filepath.Clean(fmt.Sprintf("%s/autoload/command.go", dir))
	autoloadContent, err := ioutil.ReadFile(autoloadFilePath)
	if err != nil {
		log.Fatal(err)
	}

	importName := fmt.Sprintf("%scommand", app.Value)
	commandName := fmt.Sprintf("%s.%s", importName, name.Value)
	autoloadContentStr := cHelper.ReplaceAllIfNotContain(string(autoloadContent), [][]string{
		{"import (", "func initCommand() {", "import (\n)\n\nfunc initCommand() {"},
		{"//TODO:ImportCommand", "import (", "import (\n\t//TODO:ImportCommand"},
		{"//TODO:InitCommand", "func initCommand() {", "func initCommand() {\n\t//TODO:InitCommand"},
		{"core/cServer", "//TODO:ImportCommand", fmt.Sprintf("//TODO:ImportCommand\n\t\"%s/core/cServer\"", project)},
		{importName, "//TODO:ImportCommand", fmt.Sprintf("//TODO:ImportCommand\n\t%s \"%s/app/%s/command\"", importName, project, app.Value)},
		{commandName, "//TODO:InitCommand", fmt.Sprintf("\t//TODO:InitCommand\n\tcServer.Inject(&%s{})", commandName)},
	})

	// 写文件
	err = ioutil.WriteFile(autoloadFilePath, []byte(autoloadContentStr), 0755)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("autoload success:", autoloadFilePath)
}

func (i *AppCommand) Middleware(app cCommand.Option, name cCommand.Option) {
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

	// 读取文件
	content, err := asset.Asset("resource/template/app/app/middleware/TestMiddleware.go.tpl")
	if err != nil {
		log.Fatal(err)
	}

	contentStr := cHelper.ReplaceAllFromMap(string(content), map[string]string{
		"##PROJECT##":    project,
		"##MIDDLEWARE##": name.Value,
	})

	// 写文件
	targetFilePath := filepath.Clean(fmt.Sprintf("%s/app/%s/middleware/%s.go", dir, app.Value, name.Value))
	if cHelper.IsExistsPath(targetFilePath) {
		log.Fatalf("file is exists: %s", targetFilePath)
	}
	err = ioutil.WriteFile(targetFilePath, []byte(contentStr), 0755)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("create success:", targetFilePath)

	// 更新autoload
	autoloadFilePath := filepath.Clean(fmt.Sprintf("%s/autoload/middleware.go", dir))
	autoloadContent, err := ioutil.ReadFile(autoloadFilePath)
	if err != nil {
		log.Fatal(err)
	}

	importName := fmt.Sprintf("%smiddleware", app.Value)
	middlewareName := fmt.Sprintf("%s.%s", importName, name.Value)
	autoloadContentStr := cHelper.ReplaceAllIfNotContain(string(autoloadContent), [][]string{
		{"import (", "func initMiddleware() {", "import (\n)\n\nfunc initMiddleware() {"},
		{"//TODO:ImportMiddleware", "import (", "import (\n\t//TODO:ImportMiddleware"},
		{"//TODO:InitMiddleware", "func initMiddleware() {", "func initMiddleware() {\n\t//TODO:InitMiddleware"},
		{"core/cServer", "//TODO:ImportMiddleware", fmt.Sprintf("//TODO:ImportMiddleware\n\t\"%s/core/cServer\"", project)},
		{importName, "//TODO:ImportMiddleware", fmt.Sprintf("//TODO:ImportMiddleware\n\t%s \"%s/app/%s/middleware\"", importName, project, app.Value)},
		{fmt.Sprintf("%s.%s", importName, name.Value), "//TODO:InitMiddleware", fmt.Sprintf("//TODO:InitMiddleware\n\tcServer.Inject(&%s{})", middlewareName)},
	})

	// 写文件
	err = ioutil.WriteFile(autoloadFilePath, []byte(autoloadContentStr), 0755)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("autoload success:", autoloadFilePath)
}

func (i *AppCommand) Orm(app cCommand.Option, name cCommand.Option) {}

func protoImportFiles(protoPath, name string) []string {
	var res []string

	importFilePath := filepath.Clean(fmt.Sprintf("%s/%s", protoPath, name))
	// 读取文件
	content, _ := ioutil.ReadFile(importFilePath)

	importReg := regexp.MustCompile(`import +"([0-9a-zA-Z\/]+\.proto)"`)
	importMatches := importReg.FindAllStringSubmatch(string(content), -1)

	for _, v := range importMatches {
		if v[1] != "" {
			tempRes := protoImportFiles(protoPath, v[1])
			res = append(res, v[1])
			res = append(res, tempRes...)
		}
	}

	return res
}

func protoImportPackage(protoPath, name, project string) []string {
	var res []string

	importFilePath := filepath.Clean(fmt.Sprintf("%s/%s", protoPath, name))
	// 读取文件
	content, _ := ioutil.ReadFile(importFilePath)

	importReg := regexp.MustCompile(`import +"([0-9a-zA-Z\/]+\.proto)"`)
	importMatches := importReg.FindAllStringSubmatch(string(content), -1)

	subImportReg := regexp.MustCompile(`option +go_package *= *"([0-9a-zA-Z\/\.]+)"`)
	for _, v := range importMatches {
		if v[1] != "" {
			subImportFilePath := filepath.Clean(fmt.Sprintf("%s/%s", protoPath, v[1]))
			subContent, _ := ioutil.ReadFile(subImportFilePath)
			subImportMatches := subImportReg.FindAllStringSubmatch(string(subContent), -1)
			for _, sv := range subImportMatches {
				packageName := strings.Replace(sv[1], "../", project+"/base/", 1)
				res = append(res, packageName)
			}
		}
	}

	return res
}
