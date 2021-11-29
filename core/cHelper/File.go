package cHelper

import "os"

// IsExistsPath 判断文件是否存在
func IsExistsPath(path string) bool {
	_, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		return true
	}

	return false
}

