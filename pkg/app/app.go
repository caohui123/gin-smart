package app

import (
	"flag"
	"fmt"
	"os"
)

const (
	TimeFormatFullDate = "2006-01-02 15:04:05" // 常规类型
	EnvLocal           = "local"
	EnvDev             = "dev"
	EnvTest            = "test"
	EnvProduction      = "production"
)

type TriggerIF interface {
	Do()
}

// 触发
func Trigger(tg TriggerIF) {
	go tg.Do()
}

func PrintConsole(value interface{}) {
	var content string
	if err, ok := value.(error); ok {
		content = err.Error()
	}
	if content == "" {
		content = fmt.Sprintf("%v", value)
	}
	fmt.Println("---console---", content)
}

// 控制台参数
type StarArgs struct {
	config string
}

//
func GetStartArgs() StarArgs {
	// 配置文件路径，取命令行config参数作为路径
	configFile := getConfigFile()
	cmdArgsConfig := flag.String("config", configFile, "config file path, default: "+configFile)
	flag.Parse()
	if cmdArgsConfig != nil {
		configFile = *cmdArgsConfig
	}
	return StarArgs{
		config: configFile,
	}
}

func BootPath() string {
	str, _ := os.Getwd()
	return str
}

func mustBootInProjectRoot() {
	file := BootPath() + "/go.mod"
	if ok, err := IsPathExists(file); err != nil {
		panic("check go.mod exists err: " + err.Error())
	} else if !ok {
		panic(file + " not found, please boot in project root dir")
	}
}
