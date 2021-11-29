package cCommand

// listenExit 监听退出信号
func listenExit() {
	go func() {
		for {
			select {
			case exit := <-state.ExitChannel:
				if exit {
					state.Enable = false
					StopCron()
					// 通知系统退出完成
					if state.AppExitChannel != nil {
						state.AppExitChannel <- true
					}
				}
			}
		}
	}()
}

// SetSysExitChannel 设置系统退出信号通道
func SetSysExitChannel(c chan bool) {
	state.AppExitChannel = c
}

// Exit 发出退出信号
func Exit() {
	go func() {
		state.ExitChannel <- true
	}()
}
