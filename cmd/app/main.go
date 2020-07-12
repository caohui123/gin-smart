package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/param"
	"github.com/jangozw/gin-smart/pkg/app"
	"github.com/jangozw/gin-smart/pkg/util"
	"github.com/jangozw/gin-smart/route"
	"github.com/urfave/cli/v2"
)

var (
	BuildVersion string // 编译的app版本
	BuildAt      string // 编译时间
)

func main() {
	if err := stack().Run(os.Args); err != nil {
		panic(err)
	}
}

func stack() *cli.App {
	return &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        param.ArgConfig,
				Value:       param.ArgConfigFilename,
				Destination: &app.ConfigFile,
			},
		},
		Name:    param.AppName,
		Usage:   fmt.Sprintf("./%s -%s=%s", param.AppName, param.ArgConfig, param.ArgConfigFilename),
		Version: fmt.Sprintf("%s-%s-%s-%s-%s", runtime.GOOS, runtime.GOARCH, BuildVersion, BuildAt, util.Now()),
		Action:  action,
		Before:  before,
	}
}

// 初始化加载服务
func before(c *cli.Context) error {
	app.LoadServices()
	app.BuildInfo = c.App.Version
	if app.Logger != nil {
		app.Logger.Infof("%s init successfully! loaded services: %s", param.AppName, app.Loaded)
	}
	return nil
}

// 注册路由
func action(c *cli.Context) error {
	// 注册路由
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	route.Register(engine)
	if err := engine.Run(fmt.Sprintf(":%d", app.Cfg.General.ApiPort)); err != nil {
		return err
	}
	return nil
}
