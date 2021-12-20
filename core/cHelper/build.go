package cHelper

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func GetModName() (string, error) {
	// 获取创建项目的路径
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 读取project name
	f := dir + "/go.mod"
	mod, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}

	modReg := regexp.MustCompile(`module ([\w\.\-\_]+)`)
	res := modReg.FindAllStringSubmatch(string(mod), -1)

	if len(res) == 1 && len(res[0]) == 2 {
		return res[0][1], nil
	}

	return "", errors.New("mod name not foound")
}
