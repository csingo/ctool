package cHelper

import "os"

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
