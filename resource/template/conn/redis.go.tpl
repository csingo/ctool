package conn

import (
	"fmt"
	"##PROJECT##/config/vars"
	"github.com/gomodule/redigo/redis"
)

func Redis() redis.Conn {
	dbconf := vars.Redis
	host := dbconf.Host
	port := dbconf.Port
	auth := dbconf.Auth
	db := dbconf.Db

	var options []redis.DialOption
	if auth != "" {
		options = append(options, redis.DialPassword(auth))
	}
	if db > 0 {
		options = append(options, redis.DialDatabase(db))
	}

	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port), options...)

	if err != nil {
		panic("redis connection err: " + err.Error())
	}

	return conn
}
