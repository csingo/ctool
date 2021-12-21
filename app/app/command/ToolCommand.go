package command

import (
	"archive/zip"
	"fmt"
	"github.com/csingo/ctool/config/vars"
	"github.com/csingo/ctool/core/cCommand"
	"github.com/csingo/ctool/core/cHelper"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type ToolCommand struct{}

func (i *ToolCommand) Help() *cCommand.CommandHelpDoc {
	return &cCommand.CommandHelpDoc{
		CommandDesc: cCommand.CommandDesc{Name: "tool", Desc: "工具"},
		MethodDesc: []cCommand.MethodDesc{
			{Name: "init", Desc: "初始化", Options: []string{}},
			{Name: "update", Desc: "升级", Options: []string{}},
		},
		OptionDesc: []cCommand.OptionDesc{},
	}
}

func (i *ToolCommand) Init() {
	var err error
	// 当前工作目录
	thisPath, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	// GOPATH
	gopath := cHelper.GetGOENV("GOPATH")
	envOS := cHelper.GetGOENV("GOOS")
	// 下载 protoc
	downloadFile := filepath.Clean(fmt.Sprintf("%s/%s", thisPath, "protoc.zip"))
	var downloadUrl string
	var protofile string
	switch envOS {
	case "linux":
		downloadUrl = vars.Tool.ProtocDownload.Linux
		protofile = gopath + "/bin/protoc"
	case "windows":
		downloadUrl = vars.Tool.ProtocDownload.Win
		protofile = gopath + "/bin/protoc.exe"
	case "darwin":
		downloadUrl = vars.Tool.ProtocDownload.Mac
		protofile = gopath + "/bin/protoc"
	}
	log.Println("download protoc.zip ...", downloadFile)
	downloadRsp, err := http.Get(downloadUrl)
	if err != nil {
		log.Fatalln(err)
	}
	f, err := os.Create(downloadFile)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = io.Copy(f, downloadRsp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// 解压 protoc
	log.Println("unzip protoc.zip ...", downloadFile)
	zipReader, err := zip.OpenReader(downloadFile)
	if err != nil {
		log.Fatalln(err)
	}
	for _, zf := range zipReader.File {
		if zf.Name == "bin/protoc.exe" || zf.Name == "bin/protoc" {
			inFile, err := zf.Open()
			if err != nil {
				log.Fatalln(err)
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(protofile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zf.Mode())
			if err != nil {
				log.Fatalln(err)
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				log.Fatalln(err)
			}
			break
		}
	}
	// 安装 proto-gen-go
	log.Println("install protoc-gen-go ...")
	cmd := exec.Command("go", "install", vars.Tool.ProtoGenGoPackage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
	// 安装 go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	log.Println("install protoc-gen-go-grpc ...")
	cmd = exec.Command("go", "install", vars.Tool.ProtoGenGoGrpcPackage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		f.Close()
		zipReader.Close()
		// 删除下载的源文件
		log.Println("clean protoc.zip ...")
		err = os.Remove(downloadFile)
		if err != nil {
			log.Fatalln(err)
		}
	}()
}

func (i *ToolCommand) Update() {
	modcache := cHelper.GetGOENV("GOMODCACHE")
	ctoolPath := filepath.Clean(fmt.Sprintf("%s/gitee.com/csingo", modcache))
	log.Println("clean mod cache", ctoolPath)

	// 清理mod缓存
	err := os.RemoveAll(ctoolPath)
	if err != nil {
		log.Fatalln(err)
	}

	gopath := cHelper.GetGOENV("GOPATH")
	goos := cHelper.GetGOENV("GOOS")
	var ctoolFilePath, ctoolDownloadUrl string
	switch goos {
	case "windows":
		ctoolFilePath = fmt.Sprintf("%s/%s", gopath, "bin/ctool.exe")
		ctoolDownloadUrl = "https://raw.githubusercontent.com/csingo/ctool/raw/master/bin/win-ctool.exe"
	case "linux":
		ctoolFilePath = fmt.Sprintf("%s/%s", gopath, "bin/ctool")
		ctoolDownloadUrl = "https://raw.githubusercontent.com/csingo/ctool/raw/master/bin/linux-ctool"
	case "darwin":
		ctoolFilePath = fmt.Sprintf("%s/%s", gopath, "bin/ctool")
		ctoolDownloadUrl = "https://raw.githubusercontent.com/csingo/ctool/raw/master/bin/mac-ctool"
	}

	// 下载工具
	log.Println("download ctool ...", ctoolDownloadUrl)
	downloadRsp, err := http.Get(ctoolDownloadUrl)
	if err != nil {
		log.Fatalln(err)
	}

	// 删除工具
	log.Println("uninstall ctool ...", ctoolFilePath)
	err = os.Remove(ctoolFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	// 安装工具
	log.Println("install ctool ...", ctoolFilePath)
	f, err := os.Create(ctoolFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = io.Copy(f, downloadRsp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// 授权
	log.Println("chmod ctool ...", ctoolFilePath)
	err = os.Chmod(ctoolFilePath, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("finish")
}


