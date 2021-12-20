package cHelper

import (
	"os"
	"os/exec"
	"strings"
)

func EnvToInt(name string, value int) int {
	env := os.Getenv(name)
	if env == "" {
		return value
	}

	return ToInt(env)
}

func EnvToBool(name string, value bool) bool {
	env := os.Getenv(name)
	if env == "" {
		return value
	}

	return ToBool(env)
}

func EnvToString(name string, value string) string {
	env := os.Getenv(name)
	if env == "" {
		return value
	}

	return env
}

func GetGOENV(name string) string {
	cmd := exec.Command("go", "env", name)
	env, _ := cmd.Output()
	cmd.Run()

	res := string(env)

	res = strings.Trim(res, "\r")
	res = strings.Trim(res, "\n")

	return res
}
