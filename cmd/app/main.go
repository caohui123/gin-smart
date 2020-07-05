package main

import (
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
		fmt.Println("Exit! app setup runner failed: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("init main ok, --current env:", app.CurrentEnv(), "--boot at:", app.BootPath())
}

func main() {
	before()
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	// 注册路由
	route.Register(engine)
	if err := engine.Run(fmt.Sprintf(":%d", app.Runner.Cfg.Server.Listen)); err != nil {
		panic(err)
	}
}
