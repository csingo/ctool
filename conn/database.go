package conn

import (
	"github.com/csingo/ctool/config/vars"
	"github.com/csingo/ctool/core/cHelper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type BaseModel struct{}

func (i BaseModel) Db() *gorm.DB {
	return connectToDatabase()
}

func (i BaseModel) Model(value interface{}) *gorm.DB {
	c := connectToDatabase()

	return c.Model(value)
}

func (i BaseModel) Table(name string, args ...interface{}) *gorm.DB {
	c := connectToDatabase()

	return c.Table(name, args...)
}

var db *gorm.DB = nil

func connectToDatabase() *gorm.DB {
	if db != nil {
		return db
	}
	var err error
	dbconf := vars.Database
	//driver := dbconf.Driver
	channels := dbconf.Channels

	if channels == nil || channels.Default == nil {
		panic("database config driver is not set")
	}

	// TODO: driver 未起作用
	driverConf := channels.Default
	var dsn = driverConf.Dsn
	if dsn == "" {
		dsn = driverConf.Username + ":" +
			driverConf.Password + "@tcp(" +
			driverConf.Host + ":" +
			cHelper.ToString(driverConf.Port) + ")/" +
			driverConf.Database + "?charset=" +
			driverConf.Charset + "&parseTime=" +
			driverConf.Parsetime + "&loc=" +
			driverConf.Loc
	}
	conn := mysql.Open(dsn)
	db, err = gorm.Open(conn, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	if err != nil {
		panic("database connection err: " + err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("database connection pool err: " + err.Error())
	}
	sqlDB.SetMaxIdleConns(driverConf.MaxIgleConns)
	sqlDB.SetMaxOpenConns(driverConf.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(driverConf.ConnMaxLifetime))

	return db
}
