package vars

import (
	"gitee.com/csingo/ctool/config/typs"
	"gitee.com/csingo/ctool/core/cHelper"
)

// PanicLevel 0
// FatalLevel
// ErrorLevel
// WarnLevel
// InfoLevel
// DebugLevel
// TraceLevel

var Log = &typs.LogConf{
	Level:  cHelper.EnvToInt("Log_LEVEL", 7),
	Output: cHelper.EnvToString("Log_OUTPUT", "/dev/stdout"),
}
