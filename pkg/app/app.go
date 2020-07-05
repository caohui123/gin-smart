package app

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	TimeFormatFullDate = "2006-01-02 15:04:05" // 常规类型
	EnvLocal           = "local"
	EnvDev             = "dev"
	EnvTest            = "test"
	EnvProduction      = "production"
)
var Setting SettingFields


type paramsCheck interface {
	Check() error
}
type TriggerIF interface {
	Do()
}

// 触发
func Trigger(tg TriggerIF) {
	go tg.Do()
}

// 控制台参数
type BootArgs struct {
	Config string
}

//
func GetBootArgs() BootArgs {
	// 配置文件路径，取命令行config参数作为路径
	configFile := getConfigFile()
	cmdArgsConfig := flag.String("config", configFile, "config file path, default: "+configFile)
	flag.Parse()
	if cmdArgsConfig != nil {
		configFile = *cmdArgsConfig
	}
	return BootArgs{
		Config: configFile,
	}
}

func BootPath() string {
	str, _ := os.Getwd()
	return str
}

type SettingFields struct {
	ErrCodeMap     map[int]string
	ErrCodeSuccess int
	BuildAt        string
	StartAt        time.Time
	BuildVersion   string
	BootArgs       BootArgs
}


func getCodeMsg(code int) string {
	if v, ok := Setting.ErrCodeMap[code]; ok {
		return v
	}
	return fmt.Sprintf("unknown code:%d", code)
}

// 环境
func IsEnvLocal() bool {
	return CurrentEnv() == EnvLocal
}

func IsEnvDev() bool {
	return CurrentEnv() == EnvDev
}

func IsEnvTest() bool {
	return CurrentEnv() == EnvTest
}

func IsEnvProduction() bool {
	return CurrentEnv() == EnvProduction
}

func CurrentEnv() string {
	if Runner != nil {
		return Runner.Cfg.Env.Env
	}
	return ""
}
