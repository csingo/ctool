package cConfig

import (
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"time"
)

type ConfigInterface interface {
	ConfigName() string
}

type container struct {
	path           string
	instances      map[string]interface{}
	versions       map[string]string
	updateInterval int
	updateChannel  chan string
	updateTickers  map[string]*time.Ticker
	updateFunc     map[string][]HandleFunc
	configCenter   *center
}

type center struct {
	enable      bool
	driver      string
	nacosClient config_client.IConfigClient
	listeners   listeners
	listenNames names
}

type listeners map[string]string
type names map[string]string

type HandleFunc func()

var configContainer = &container{
	path:           "",
	instances:      make(map[string]interface{}),
	versions:       make(map[string]string),
	updateInterval: 1,
	updateChannel:  make(chan string),
	updateTickers:  make(map[string]*time.Ticker),
	updateFunc:     make(map[string][]HandleFunc),
	configCenter: &center{
		enable:      false,
		nacosClient: nil,
		listeners:   make(listeners),
		listenNames: make(map[string]string),
	},
}
