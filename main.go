package main

import (
	"log"
	"reflect"
)

type test struct {

}

func (receiver test) Name(a string) {
	log.Println(a)
}

var t = &test{}

func main() {
	caller := reflect.ValueOf(t).Elem().MethodByName("Name")
	param := caller.Type().In(0)
	a := reflect.New(param)

	log.Printf("%#V", a)

	//autoload.Init()
	//cServer.Init()
	//
	//cServer.Start()
}
