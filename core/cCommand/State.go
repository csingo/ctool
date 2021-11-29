package cCommand

// listenExit 监听退出信号
func listenExit() {
	go func() {
		for {
			select {
			case exit := <-state.exitChannel:
				if exit {
					state.enable = false
					StopCron()
					// 通知系统退出完成
					if state.appExitChannel != nil {
						state.appExitChannel <- true
					}
				}
			}
		}
	}()
}

// SetSysExitChannel 设置系统退出信号通道
func SetSysExitChannel(c chan bool) {
	state.appExitChannel = c
}

// Exit 发出退出信号
func Exit() {
	go func() {
		state.exitChannel <- true
	}()
}
