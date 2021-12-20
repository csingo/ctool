package cLog

import (
	"github.com/csingo/ctool/core/cConfig"
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Load() {
	// 设置将日志输出到标准输出
	log.SetOutput(os.Stdout)
	// 设置日志格式为json格式和时间格式
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		DataKey:         "context",
	})
	// 设置日志级别
	level, err := cConfig.GetConf("LogConf.Level")
	if err != nil {
		level = 7
	}
	log.SetLevel(logrus.Level(uint32(level.(int))))

}

func Trace(msg string, parmas map[string]interface{}) {
	WithFields(parmas).Trace(msg)
}

func Debug(msg string, parmas map[string]interface{}) {
	WithFields(parmas).Debug(msg)
}

func Info(msg string, parmas map[string]interface{}) {
	WithFields(parmas).Info(msg)
}

func Warn(msg string, parmas map[string]interface{}) {
	WithFields(parmas).Warn(msg)
}

func Error(msg string, parmas map[string]interface{}) {
	WithFields(parmas).Error(msg)
}

func Fatal(msg string, parmas map[string]interface{}) {
	WithFields(parmas).Fatal(msg)
}

func Panic(msg string, parmas map[string]interface{}) {
	WithFields(parmas).Panic(msg)
}

func WithFields(parmas map[string]interface{}) *logrus.Entry {
	// todo 这里后续设置trace_id
	return log.WithFields(parmas)
}
