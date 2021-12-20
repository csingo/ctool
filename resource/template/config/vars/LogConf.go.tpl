package vars

import (
	"##PROJECT##/config/typs"
	"##PROJECT##/core/cHelper"
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
