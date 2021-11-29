package cServer

import (
	"net/http"
	"os"
	"time"
)

type container struct {
	Instances map[string]interface{}
}

var ServerContainer = &container{
	Instances: make(map[string]interface{}),
}

type application struct {
	Signal        chan os.Signal
	Exit          chan bool
	ExitCounter   uint
	EnableHTTP    chan bool
	EnableTCP     chan bool
	EnableCommand chan bool
	EnableConfig  chan bool
}

var App = &application{
	Signal:        make(chan os.Signal),
	Exit:          make(chan bool),
	ExitCounter:   0,
	EnableHTTP:    make(chan bool),
	EnableTCP:     make(chan bool),
	EnableCommand: make(chan bool),
	EnableConfig:  make(chan bool),
}

type ServerState struct {
	Enable          bool
	ExitChannel     chan bool
	HTTPConnCounter int
	HTTPServer      *http.Server
	Timeout         time.Duration
}

var state = &ServerState{
	Enable:          true,
	ExitChannel:     make(chan bool),
	HTTPConnCounter: 0,
	HTTPServer:      nil,
	Timeout:         5 * time.Second,
}


