package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/api/sample"
	"github.com/jangozw/gin-smart/middleware"
	"github.com/jangozw/gin-smart/pkg/app"
)

var warp = app.Warp

// 注册路由
func Register(router *gin.Engine) {
	sampleGroup := router.Group("/sample", warp(middleware.LogRequest), warp(middleware.Header))
	homeGroup := router.Group("/", warp(middleware.LogRequest), warp(middleware.Header))
	v1Group := router.Group("/v1", warp(middleware.LogRequest), warp(middleware.Header))
	v1LoginGroup := router.Group("/v1", warp(middleware.LogRequest), warp(middleware.Header), warp(middleware.CheckToken))

	sampleGroup.GET("", warp(func(context *gin.Context) {
		ctx := app.Ctx(context)
		ctx.MustLoginUser()
		ctx.Success("Sample here")
	}))

	sampleGroup.GET("sample", warp(sample.Sample))

	sampleGroup.GET("/list", warp(sample.List))

	homeGroup.GET(`/home`, warp(func(context *gin.Context) {
		ctx := app.Ctx(context)
		user := ctx.MustLoginUser()
		// panic("aa")
		ctx.Success(user)
	}))

	v1Group.GET("/test", warp(func(c *gin.Context) {
		app.Ctx(c).Success("ok v1 test")
	}))

	v1LoginGroup.GET(`/config`, warp(func(c *gin.Context) {
		app.Ctx(c).Success("ok")
	}))
}
