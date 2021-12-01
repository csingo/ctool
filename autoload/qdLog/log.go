package qdLog

import (
	"os"

	"gitee.com/csingo/ctool/config/vars"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func InitLog() {
	// 设置将日志输出到标准输出
	log.SetOutput(os.Stdout)
	// 设置日志格式为json格式和时间格式
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		DataKey:         "context",
	})
	// 设置日志级别
	log.SetLevel(logrus.Level(uint32(vars.Log.Level)))

}

func Trace(msg string, parmas map[string]interface{}) {
	WithFields(parmas).Info(msg)
}

func Debug(msg string, parmas map[string]interface{}) {
	WithFields(parmas).Info(msg)
}

func Info(msg string, parmas map[string]interface{}) {
	WithFields(parmas).Info(msg)
}

func Warn(msg string, parmas map[string]interface{}) {
	WithFields(parmas).Info(msg)
}

func Error(msg string, parmas map[string]interface{}) {
	WithFields(parmas).Info(msg)
}

func Fatal(msg string, parmas map[string]interface{}) {
	WithFields(parmas).Info(msg)
}

func Panic(msg string, parmas map[string]interface{}) {
	WithFields(parmas).Info(msg)
}

func WithFields(parmas map[string]interface{}) *logrus.Entry {
	// todo 这里后续设置trace_id
	return log.WithFields(parmas)
}
