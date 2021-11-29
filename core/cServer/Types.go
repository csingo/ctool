package cServer

import (
	"net/http"
	"os"
	"time"
)

type container struct {
	instances map[string]interface{}
}

var serverContainer = &container{
	instances: make(map[string]interface{}),
}

type application struct {
	signal        chan os.Signal
	exit          chan bool
	exitCounter   uint
	enableHTTP    chan bool
	enableTCP     chan bool
	enableCommand chan bool
	enableConfig  chan bool
}

var app = &application{
	signal:        make(chan os.Signal),
	exit:          make(chan bool),
	exitCounter:   0,
	enableHTTP:    make(chan bool),
	enableTCP:     make(chan bool),
	enableCommand: make(chan bool),
	enableConfig:  make(chan bool),
}

type serverState struct {
	enable          bool
	exitChannel     chan bool
	httpConnCounter int
	httpServer      *http.Server
	timeout         time.Duration
}

var state = &serverState{
	enable:          true,
	exitChannel:     make(chan bool),
	httpConnCounter: 0,
	httpServer:      nil,
	timeout:         5 * time.Second,
}


