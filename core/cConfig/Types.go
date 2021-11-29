package cConfig

import (
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"time"
)

type ConfigInterface interface {
	ConfigName() string
}

type container struct {
	Path           string
	Instances      map[string]interface{}
	Versions       map[string]string
	UpdateInterval int
	UpdateChannel  chan string
	UpdateTickers  map[string]*time.Ticker
	UpdateFunc     map[string][]HandleFunc
	ConfigCenter   *center
}

type center struct {
	Enable      bool
	Driver      string
	NacosClient config_client.IConfigClient
	Listeners   listeners
}

type listeners map[string]string

type HandleFunc func()

var ConfigContainer = &container{
	Path:           "",
	Instances:      make(map[string]interface{}),
	Versions:       make(map[string]string),
	UpdateInterval: 1,
	UpdateChannel:  make(chan string),
	UpdateTickers:  make(map[string]*time.Ticker),
	UpdateFunc:     make(map[string][]HandleFunc),
	ConfigCenter: &center{
		Enable:      false,
		NacosClient: nil,
		Listeners:   make(listeners),
	},
}
