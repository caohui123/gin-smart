package main

import (
	"encoding/json"
	"fmt"
	"github.com/jangozw/gin-smart/errcode"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/route"

	"github.com/jangozw/gin-smart/pkg/app"
)

var (
	BuildVersion string // 编译的app版本
	BuildAt      string // 编译时间
)

func before() {
	app.Setting = app.SettingFields{
		BootArgs:       app.GetBootArgs(),
		BuildVersion:   BuildVersion,
		BuildAt:        BuildAt,
		StartAt:        time.Now(),
		ErrCodeMap:     errcode.CodeMap(),
		ErrCodeSuccess: errcode.Success,
	}
	// 注册运行依赖的资源,db,redis等
	if err := app.NewRunner(); err != nil {
		fmt.Println("Exit! setup app.Runner failed: ", err.Error())
		os.Exit(1)
	}
	printInitInfo()
}

func main() {
	before()
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	// 注册路由
	route.Register(engine)
	if err := engine.Run(fmt.Sprintf(":%d", app.Runner.Cfg.Server.Listen)); err != nil {
		fmt.Println("Exit! setup web server failed: ", err.Error())
		os.Exit(1)
	}
}

func printInitInfo()  {
	var info = map[string]string{
		"AppEnv":app.CurrentEnv(),
		"Build": app.Setting.BuildVersion+`@`+app.Setting.BuildAt,
		"Config":app.GetConfigFile(),
		"LogDir":app.Runner.Cfg.Log.LogDir,
		"BootAt":app.BootPath(),
	}
	by,_:= json.Marshal(info)
	fmt.Println("App init completely! ", string(by))
}
