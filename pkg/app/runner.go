package app

import (
	"time"

	"github.com/jangozw/gin-smart/pkg/lib"
)

// 默认是启动根目录下的 config.ini ， 启动时候可以-config 参数指定任意路径
var ConfigPath string

// 运行的资源实例
var Runner *services

// 主程序启动时的信息
var MainStartInfo StartInformation

// 运行加载的实例
func NewRunner() (*services, error) {
	conf, err := lib.NewConfig(ConfigPath)
	if err != nil {
		return nil, err
	}
	db, err := lib.NewDb(conf)
	if err != nil {
		return nil, err
	}
	redis, err := lib.NewRedis(conf)
	if err != nil {
		return nil, err
	}
	logger, err := lib.NewLogger(conf)
	if err != nil {
		return nil, err
	}
	Runner = &services{
		Conf:  conf,
		Db:    db,
		Redis: redis,
		Log:   logger,
	}
	return Runner, nil
}

// 运行的资源
type services struct {
	Db    *lib.DbMysql
	Redis *lib.RedisPool
	Log   *lib.Logger
	Conf  *lib.Config
}

// 启动信息
type StartInformation struct {
	BuildAt      string
	BuildVersion string
	StartAt      time.Time
	Root         string
	ConfigPath   string
}
