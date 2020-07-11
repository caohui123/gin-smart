package app

import (
	"errors"
	"fmt"
	"github.com/jangozw/gin-smart/config"
	"github.com/jangozw/gin-smart/pkg/lib"
	"github.com/jangozw/gin-smart/pkg/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"strings"
)

// 配置文件路径，启动参数指定
var ConfigFile string

// 解析的配置文件
var Cfg *config.Config

// 数据库
var Db *gorm.DB

// redis
var Redis *lib.Redis

// 日志
var Log *logrus.Logger

const (
	LogService   = `Log`
	RedisService = `Redis`
	DbService    = `Db`
)

// app 需要自动加载的服务配置，用不到的可以注释掉
type loadServiceMap map[string]func() error
var serviceMap = loadServiceMap{
	LogService:   LoadLog,
	RedisService: LoadRedis,
	DbService:    LoadDb,
}

// app 包初始化，加载服务失败要panic
func LoadServices(codes ...string) {
	LoadCfg()
	if len(codes) == 0 {
		codes = serviceMap.keys()
	}
	for key, load := range serviceMap {
		if InStringSlice(key, codes) {
			if err := load(); err != nil {
				panic("app加载服务失败:" + err.Error())
			}
		}
	}
	fmt.Println("--app加载服务完成, loaded:", strings.Join(codes, `,`))
}

// loadCfg 从启动参数或者项目目录中查找并加载配置文件
func LoadCfg() {
	if Cfg != nil {
		return
	}
	// ConfigFile 在启动时候赋值
	if ConfigFile == "" {
		// 配置文件的启动参数名称 eg: -config=xx.ini
		var configFlag = "config"
		// 配置文件相对于运行目录的路径
		var filename = "config.ini"

		var err = errors.New(fmt.Sprintf("查找配置文件%s出错: %s", filename, "文件不存在"))
		f, _ := util.FindConfigFile(filename, configFlag)
		if f == "" {
			panic(err)
		}
		ok, _ := util.IsPathExists(f)
		if !ok {
			panic(err)
		}
		ConfigFile = f
	}
	Cfg = &config.Config{}
	if err := util.ParseIni(ConfigFile, Cfg); err != nil {
		panic(fmt.Sprintf("解析配置文件%s出错: %s", ConfigFile, err.Error()))
	}
}

func LoadDb() (err error) {
	if Db != nil {
		return nil
	}
	LoadCfg()
	cfgDb := lib.CfgDatabase{
		Schema:   Cfg.Database.Schema,
		Host:     Cfg.Database.Host,
		User:     Cfg.Database.User,
		Password: Cfg.Database.Password,
		Database: Cfg.Database.Database,
	}
	Db, err = lib.NewDb(cfgDb)
	return
}

func LoadRedis() (err error) {
	if Redis != nil {
		return nil
	}
	LoadCfg()
	cfgRedis := lib.CfgRedis{
		Host:     Cfg.Redis.Host,
		Password: Cfg.Redis.Password,
		DbNum:    Cfg.Redis.DbNum,
	}
	Redis, err = lib.NewRedis(cfgRedis)
	return
}

func LoadLog() (err error) {
	if Log != nil {
		return nil
	}
	LoadCfg()
	Log, err = lib.NewLogger(Cfg.General.LogDir, "app")
	return
}



func InStringSlice(key string, src []string) bool {
	for _, v := range src {
		if v == key {
			return true
		}
	}
	return false
}

func (m loadServiceMap) keys() []string {
	keys := make([]string, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}
