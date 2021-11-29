package cConfig

import (
	"os"
	"path/filepath"
	"strings"
)

func initLocalFile() {
	var path string
	argsLen := len(os.Args)
	if argsLen < 2 {
		path, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	} else {
		path = os.Args[1]
	}
	path = strings.TrimRight(path, "/")
	path = strings.TrimRight(path, "\\")
	configContainer.path = path
}
