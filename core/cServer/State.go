package cServer

// listenExit 监听退出信号
func listenExit() {
	go func() {
		for {
			select {
			case exit := <-state.ExitChannel:
				if exit {
					state.Enable = false
					// 退出Http服务
					StopHTTP()
					// 通知系统退出完成
					App.Exit <- true
				}
			}
		}
	}()
}

// Exit 发出退出信号
func Exit() {
	go func() {
		state.ExitChannel <- true
	}()
}
