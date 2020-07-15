package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/api/sample"
	"github.com/jangozw/gin-smart/config"
	"github.com/jangozw/gin-smart/erron"
	"github.com/jangozw/gin-smart/middleware"
	"github.com/jangozw/gin-smart/pkg/app"
)

// 注册路由
func Register(router *app.Engine) {
	// 路由组
	var (
		sampleGroup  = router.Group("/sample", middleware.LogRequest, middleware.Header)
		homeGroup    = router.Group("/", middleware.LogRequest, middleware.Header)
		v1Group      = router.Group("/v1", middleware.LogRequest, middleware.Header)
		v1LoginGroup = router.Group("/v1", middleware.LogRequest, middleware.Header, middleware.NeedLogin)
	)
	{
		// 定义路由
		homeGroup.GET("", func(c *gin.Context) (data interface{}, err erron.E) {
			return "welcome", nil
		})
	}

	// sample 组
	{
		// sampleGroup 无需登陆的路由
		sampleGroup.POST("/login", sample.Login)

		// sampleGroup 需要登陆的路由
		{
			sampleGroup.Use(middleware.NeedLogin)

			// 退出登陆
			sampleGroup.POST("logout", sample.Logout)
			sampleGroup.POST("/user", sample.AddUser)
			sampleGroup.GET("/user/list", sample.UserList)
			sampleGroup.GET("/user/detail", sample.UserDetail)
		}
	}

	// v1免登陆组
	{
		v1Group.GET("/test", func(c *gin.Context) (data interface{}, err erron.E) {
			return "test", nil
		})
	}

	// v1需要登陆组
	{
		v1LoginGroup.GET(`/config`, func(c *gin.Context) (data interface{}, err erron.E) {
			app.MustGetLoginUser(c)
			return config.GetAllStates(), nil
		})
	}
}
