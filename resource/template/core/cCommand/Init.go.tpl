package cCommand

import (
	"gitee.com/csingo/ctool/core/cConfig"
	"gitee.com/csingo/ctool/core/cHelper"
	"log"
	"reflect"
	"strings"
	"time"
)

// Inject 注入command
func Inject(instance interface{}) {
	if !IsCommand(instance) {
		return
	}

	doc := instance.(CommandInterface).Help()
	instanceName := doc.CommandDesc.Name
	instanceMethods := doc.MethodDesc
	CommandContainer.instances[instanceName] = instance

	var labels []string
	for _, method := range instanceMethods {
		labels = append(labels, instanceName+"::"+method.Name)
	}
	CommandContainer.labels = append(CommandContainer.labels, labels...)
}

// Load 初始化加载
func Load() {
	StartCron()
	StartResident()
	listenExit()
}

// IsCommand 判断是否command
func IsCommand(o interface{}) bool {
	return reflect.TypeOf(o).Implements(reflect.TypeOf((*CommandInterface)(nil)).Elem())
}

// Run 执行command
func Run(label string, options ...string) {
	if !cHelper.InArrayString(label, CommandContainer.labels) {
		log.Printf("cCommand not found: %s %v", label, options)
		return
	}
	if !state.enable {
		return
	}

	commandArr := []string{label}
	commandArr = append(commandArr, options...)
	command := strings.Join(commandArr, " ")
	labels := strings.Split(label, "::")
	commandName := labels[0]
	methodName := labels[1]

	instance := CommandContainer.instances[commandName]
	instanceDoc := instance.(CommandInterface).Help()
	methods := instanceDoc.MethodDesc
	var optionNames []string
	for _, method := range methods {
		if method.Name == methodName {
			optionNames = method.Options
		}
	}

	var optionMap = make(map[string]string)
	for _, option := range options {
		optionArr := strings.Split(option, "=")
		if len(optionArr) < 1 {
			continue
		}
		if len(optionArr) < 2 {
			optionArr = append(optionArr, "")
		}
		optionName := strings.TrimPrefix(optionArr[0], "--")
		optionValue := optionArr[1]

		optionMap[optionName] = optionValue
	}

	var commandOptions []*Option
	for _, item := range optionNames {
		if _, ok := optionMap[item]; ok {
			commandOptions = append(commandOptions, &Option{
				Name:  item,
				Value: optionMap[item],
			})
		}
	}

	argv := &CommandArgvs{
		Command: instance.(CommandInterface),
		Method:  reflect.ValueOf(instance).MethodByName(cHelper.Ucfirst(methodName)),
		Options: commandOptions,
		Name:    command,
	}

	exec(argv)
}

func exec(argv *CommandArgvs) {
	var params []reflect.Value
	for _, v := range argv.Options {
		params = append(params, reflect.ValueOf(v).Elem())
	}

	// 添加并发锁控制，生成任务uuid(纳秒+序号)， 更新任务
	state.lock.Lock()
	var uuid string
	var count int
	now := time.Now()
	t := now.Format("20060102-150405")
	ns := now.Nanosecond()
	for count = 0; count <= 999; count++ {
		uuid = cHelper.ToString(t) + "-" + cHelper.ToString(ns) + "-" + cHelper.ToString(count)
		if _, ok := state.running[uuid]; !ok {
			state.running[uuid] = argv.Name
			break
		}
	}
	state.lock.Unlock()

	defer func(uuid string) {
		if r := recover(); r != nil {
			log.Printf("cCommand exec error: %s", state.running[uuid])
		}

		// cCommand 执行完成，删除uuid
		state.lock.Lock()
		delete(state.running, uuid)
		state.lock.Unlock()
	}(uuid)

	// 执行command
	argv.Method.Call(params)
}

func StartCron() {
	// 读取定时任务配置
	cronds, err := cConfig.GetConf("CommandConf.Cronds")
	if err == nil && reflect.ValueOf(cronds).Len() > 0 {
		len := reflect.ValueOf(cronds).Len()
		for i := 0; i < len; i++ {
			cron := reflect.ValueOf(cronds).Index(i).Elem().FieldByName("Cron").Interface().(string)
			command := reflect.ValueOf(cronds).Index(i).Elem().FieldByName("Command").Interface().(string)
			options := reflect.ValueOf(cronds).Index(i).Elem().FieldByName("Options").Interface().([]string)

			var argvCrond = cron
			var argvCommand = command
			var argvOptions = options

			var cronFunc = func() {
				go Run(argvCommand, argvOptions...)
			}

			state.cron.AddFunc(argvCrond, cronFunc)
		}
	}
	// 启动定时任务
	state.cron.Start()

}

func StartResident() {
	// 读取常驻任务配置
	residents, err := cConfig.GetConf("CommandConf.Residents")
	if err == nil && reflect.ValueOf(residents).Len() > 0 {
		len := reflect.ValueOf(residents).Len()
		for i := 0; i < len; i++ {
			wait := reflect.ValueOf(residents).Index(i).Elem().FieldByName("Wait").Interface().(int)
			command := reflect.ValueOf(residents).Index(i).Elem().FieldByName("Command").Interface().(string)
			options := reflect.ValueOf(residents).Index(i).Elem().FieldByName("Options").Interface().([]string)

			var argvWait = wait
			var argvCommand = command
			var argvOptions = options

			var residentFunc = func() {
				go Run(argvCommand, argvOptions...)
			}

			time.AfterFunc(time.Duration(argvWait)*time.Second, residentFunc)
		}
	}
}

func StopCron() {
	state.cron.Stop()
}
