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

type SdkCommand struct{}

func (i *SdkCommand) Help() *cCommand.CommandHelpDoc {
	return &cCommand.CommandHelpDoc{
		CommandDesc: cCommand.CommandDesc{Name: "sdk", Desc: "服务调用sdk"},
		MethodDesc: []cCommand.MethodDesc{
			{Name: "create", Desc: "创建", Options: []string{"app", "protoPath"}},
		},
		OptionDesc: []cCommand.OptionDesc{
			{Name: "name", Desc: "名称"},
			{Name: "protoPath", Desc: "proto目录, 完整路径"},
		},
	}
}

func (i *SdkCommand) Create(app cCommand.Option, protoPath cCommand.Option) {
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

	serviceTplContent, err := ioutil.ReadFile(tplServiceFilePath)
	if err != nil {
		log.Fatal(err)
	}
	rpcTplContent, err := ioutil.ReadFile(tplRpcFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// 生成 httprpc server 和 http client
	for _, service := range services {
		servicePbFilePath := fmt.Sprintf("%s/base/%s/%s_http.pb.go", dir, app.Value, service.Name)

		var contentByte = serviceTplContent
		var content string
		content = string(contentByte)
		content = strings.ReplaceAll(content, "##SERVICE##", service.Name)
		content = strings.ReplaceAll(content, "##APP##", app.Value)

		for _, rpc := range service.Rpc {
			var subContentByte = rpcTplContent
			var subContent string
			subContent = string(subContentByte)
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

	// 执行 go mod tidy
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
