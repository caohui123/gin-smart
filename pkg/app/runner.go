package app

import (
	"github.com/jinzhu/gorm"

	"github.com/jangozw/gin-smart/pkg/lib"
)

// 运行的服务的运行者
var Runner *runner

// 运行加载的实例
func NewRunner() error {
	if Runner != nil {
		return nil
	}
	Runner = &runner{}
	// 要加载的运行服务
	runnerLoads := []func() error{
		Runner.loadCfg,
		Runner.loadLog,
		Runner.loadDb,
		Runner.loadRedis,
	}
	for _, load := range runnerLoads {
		if err := load(); err != nil {
			return err
		}
	}
	return nil
}

// 运行的资源
type runner struct {
	Db    *gorm.DB
	Redis *lib.Redis
	Log   *lib.Log
	Cfg   *Config
}

func (s *runner) loadCfg() error {
	if s.Cfg == nil {
		cfg := Config{}
		if err := lib.NewCfg(GetConfigFile(), &cfg); err != nil {
			return err
		} else if err := cfg.Check(); err != nil {
			return err
		}
		s.Cfg = &cfg
	}
	return nil
}

func (s *runner) loadDb() error {
	if err := s.loadCfg(); err != nil {
		return nil
	}
	if s.Db == nil {
		cfgDb := lib.CfgDatabase{
			Schema:   s.Cfg.Database.Schema,
			Host:     s.Cfg.Database.Host,
			User:     s.Cfg.Database.User,
			Password: s.Cfg.Database.Password,
			Database: s.Cfg.Database.Database,
		}
		if db, err := lib.NewDb(cfgDb); err != nil {
			return err
		} else {
			s.Db = db
		}
	}
	return nil
}

func (s *runner) loadRedis() error {
	if err := s.loadCfg(); err != nil {
		return nil
	}
	if s.Redis == nil {
		cfgRedis := lib.CfgRedis{
			Host:     s.Cfg.Redis.Host,
			Password: s.Cfg.Redis.Password,
			DbNum:    s.Cfg.Redis.DbNum,
		}
		if redis, err := lib.NewRedis(cfgRedis); err != nil {
			return err
		} else {
			s.Redis = redis
		}
	}
	return nil
}

func (s *runner) loadLog() error {
	if err := s.loadCfg(); err != nil {
		return nil
	}
	if s.Log == nil {
		cfgLog := lib.CfgLog{
			LogDir: s.Cfg.Log.LogDir,
		}
		if log, err := lib.NewLog(cfgLog); err != nil {
			return err
		} else {
			s.Log = log
		}
	}
	return nil
}

// 获取配置文件，启动参数指定的优先
func GetConfigFile() string {
	file := Setting.BootArgs.Config
	if file == "" {
		file = BootPath() + `/config.ini`
	}
	return file
}
