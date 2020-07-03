package lib

import "C"
import (
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // 这个不能删
	"github.com/jinzhu/gorm"
)

// use Db
type DbMysql struct {
	*gorm.DB
}

type connectConf struct {
	driver string
	args   string
}

// 初始化/conf.ini 中database区指定的数据库信息

func NewDb(appConf *Config) (*DbMysql, error) {
	section := "database"
	if v, _ := appConf.Get(section, "enable"); v == "false" {
		return nil, nil
	}
	db, err := initDatabase(appConf, section)
	if err != nil {
		return nil, err
	}
	return &DbMysql{db}, nil
}

// load database config info from /conf.ini
func initDatabase(appConf *Config, confSection string) (*gorm.DB, error) {
	c := getConnectConf(appConf, confSection)
	db, err := gorm.Open(c.driver, c.args)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("couldn't connect to database, check your connect args in config.ini, errMsg: %s", err.Error()))
	}
	// config gorm db
	db.SingularTable(true) // 全局设置表名不可以为复数形式。
	db.DB().SetMaxIdleConns(20)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(3600 * time.Second)
	// log all sql in console

	logMode, _ := appConf.GetValue("gorm", "log_mode")
	if logMode == "true" {
		db.LogMode(true)
	} else {
		db.LogMode(false)
	}

	db.Callback().Create().Replace("gorm:update_time_stamp", func(scope *gorm.Scope) {
		scope.SetColumn("CreatedAt", time.Now().Unix())
		scope.SetColumn("UpdatedAt", time.Now().Unix())
	})
	db.Callback().Update().Replace("gorm:update_time_stamp", func(scope *gorm.Scope) {
		scope.SetColumn("UpdatedAt", time.Now().Unix())
	})
	return db, nil
}

// 数据库连接配置
func getConnectConf(appConf *Config, section string) *connectConf {
	c, err := appConf.Section(section)
	if err != nil {
		panic("couldn't get config info by section:" + section)
	}
	return &connectConf{
		driver: c["schema"],
		args:   fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=%s", c["user"], c["password"], c["host"], c["dbName"], c["timezone"]),
	}
}
