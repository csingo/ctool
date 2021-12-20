package cLog

import (
	"##PROJECT##/core/cConfig"
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Load() {
	// 读取日志级别
	level, err := cConfig.GetConf("LogConf.Level")
	if err != nil {
		level = 7
	}
	// 读取日志路径
	output, err := cConfig.GetConf("LogConf.Output")
	if err != nil {
		output = "/dev/stdout"
	}

	// 设置将日志输出到标准输出
	f, err := os.OpenFile(output.(string), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		f = os.Stdout
	}
	log.SetOutput(f)
	// 设置日志级别
	log.SetLevel(logrus.Level(uint32(level.(int))))
	// 设置日志格式为json格式和时间格式
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		DataKey:         "context",
	})
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
