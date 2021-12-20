package vars

import (
	"github.com/csingo/ctool/config/typs"
	"github.com/csingo/ctool/core/cHelper"
)

// PanicLevel 0
// FatalLevel
// ErrorLevel
// WarnLevel
// InfoLevel
// DebugLevel
// TraceLevel

var Log = &typs.LogConf{
	Level: cHelper.EnvToInt("Log_LEVEL", 7),
}
