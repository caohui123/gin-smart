package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/pkg/util"
	"github.com/jangozw/gin-smart/route"
	"github.com/urfave/cli/v2"
	"os"
	"runtime"

	"github.com/jangozw/gin-smart/pkg/app"
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
				Name:        "config",
				Value:       "config.ini",
				Destination: &app.ConfigFile,
			},
		},
		Name:    "gin-smart",
		Usage:   "eg: ./gin-smart",
		Version: fmt.Sprintf("%s-%s-%s-%s-%s", runtime.GOOS, runtime.GOARCH, BuildVersion, BuildAt, util.Now()),
		Action:  action,
	}
}


func action(c *cli.Context) error {
	app.LoadServices()
	app.BuildInfo = c.App.Version
	// 注册路由
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	route.Register(engine)
	if err := engine.Run(fmt.Sprintf(":%d", app.Cfg.General.ApiPort)); err != nil {
		return err
	}
	return nil
}
