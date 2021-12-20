package cServer

// listenExit 监听退出信号
func listenExit() {
	go func() {
		for {
			select {
			case exit := <-state.exitChannel:
				if exit {
					state.enable = false
					// 退出Http服务
					StopHTTP()
					// 通知系统退出完成
					app.exit <- true
				}
			}
		}
	}()
}

// Exit 发出退出信号
func Exit() {
	go func() {
		state.exitChannel <- true
	}()
}
