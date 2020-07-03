package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.gosccap.cn/bourse/avian/errcode"
	"gitlab.gosccap.cn/bourse/avian/route"

	"gitlab.gosccap.cn/bourse/avian/pkg/app"
)

var (
	BuildVersion string // 编译的app版本
	BuildAt      string // 编译时间
)

func main() {
	before()
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	// 注册路由
	route.Register(engine)
	if port, err := app.Runner.Conf.WebPort(); err != nil {
		panic(err)
	} else if err := engine.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(err)
	}
}

func before() {
	args := readConsoleArgs()
	app.ConfigPath = args.config
	app.MainStartInfo = app.StartInformation{
		BuildVersion: BuildVersion,
		BuildAt:      BuildAt,
		Root:         bootRoot(),
		ConfigPath:   args.config,
		StartAt:      time.Now(),
	}
	// 注册运行依赖的资源,db,redis等
	runner, err := app.NewRunner()
	if err != nil {
		panic("App start failed : " + err.Error())
	}
	app.Runner = runner

	if app.IsEnvLocal() || app.IsEnvDev() {
		mustBootInProjectRoot()
	}
	app.InitSetting(app.Setting{
		ErrCodeMap: errcode.GetAllCodeMap(),
	})
	fmt.Println("init main ok, --current env:", app.Runner.Conf.Env(), "--boot root:", bootRoot())
}

// 控制台参数
type consoleArgs struct {
	config string
}

func readConsoleArgs() consoleArgs {
	// 配置文件路径，取命令行config参数作为路径
	configFile := bootRoot() + "/config.ini"
	cmdArgsConfig := flag.String("config", configFile, "config file path, default: "+configFile)
	flag.Parse()
	if cmdArgsConfig != nil {
		configFile = *cmdArgsConfig
	}
	return consoleArgs{
		config: configFile,
	}
}

func bootRoot() string {
	str, _ := os.Getwd()
	return str
}

// 本地环境检查启动，必须在项目根目录主要是用来给项目自动生成一些必要文件，其他环境在哪启动随意
func mustBootInProjectRoot() {
	file := bootRoot() + "/go.mod"
	if ok, err := isPathExists(file); err != nil {
		panic("check go.mod exists err: " + err.Error())
	} else if !ok {
		panic(file + " not found, please boot in project root dir")
	}
}

// 判断文件夹是否存在
func isPathExists(dirPath string) (bool, error) {
	_, err := os.Stat(dirPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
