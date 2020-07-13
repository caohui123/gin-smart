package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/api/sample"
	"github.com/jangozw/gin-smart/config"
	"github.com/jangozw/gin-smart/erro"
	"github.com/jangozw/gin-smart/middleware"
	"github.com/jangozw/gin-smart/pkg/app"
)

// 注册路由
func Register(engine *gin.Engine) {

	// 路由组
	var (
		router       = app.Engine(engine)
		sampleGroup  = router.Group("/sample", middleware.LogRequest, middleware.Header)
		homeGroup    = router.Group("/", middleware.LogRequest, middleware.Header)
		v1Group      = router.Group("/v1", middleware.LogRequest, middleware.Header)
		v1LoginGroup = router.Group("/v1", middleware.LogRequest, middleware.Header, middleware.CheckToken)
	)
	// 定义路由
	homeGroup.GET(``, func(c *app.Context) (data interface{}, err erro.E) {
		c.MustLoginUser()
		return "welcome", nil
	})
	sampleGroup.GET(``, func(c *app.Context) (data interface{}, err erro.E) {
		return "sample test ok", nil
	})

	sampleGroup.GET("sample", sample.Sample)

	sampleGroup.GET("/list", sample.List)

	v1Group.GET("/test", func(c *app.Context) (data interface{}, err erro.E) {
		return "test", nil
	})

	v1LoginGroup.GET(`/config`, func(c *app.Context) (data interface{}, err erro.E) {
		c.MustLoginUser()
		return config.GetAllStates(), nil
	})
}
