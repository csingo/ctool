package command

import (
	"fmt"
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

	// 执行 go mod tidy
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
