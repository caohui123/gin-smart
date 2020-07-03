package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/middleware"
	"github.com/jangozw/gin-smart/pkg/app"
	"github.com/jangozw/gin-smart/service"
)

func Register(router *gin.Engine) {
	// 用户token验证接口实例
	tokenVerifyInstance := &service.JwtUserVerify{}

	homeGroup := app.NewRouteGroup("/")
	// 参数靠前中间件的后执行
	// v0LoginGroup := app.NewRouteGroup("/v0").Use(middleware.LogToFile(), middleware.CommonMiddleware, middleware.CheckToken(tokenVerifyInstance))
	// v0FreeGroup := app.NewRouteGroup("/v0").Use(middleware.LogToFile(), middleware.CommonMiddleware)

	v1LoginGroup := app.NewRouteGroup("/v1").Use(middleware.LogToFile(), middleware.CommonMiddleware, middleware.CheckToken(tokenVerifyInstance))
	v1FreeGroup := app.NewRouteGroup("/v1").Use(middleware.LogToFile(), middleware.CommonMiddleware)
	registerRoot(homeGroup)
	// registerV0Free(v0FreeGroup)
	// registerV0Login(v0LoginGroup)

	registerV1Free(v1FreeGroup)
	registerV1Login(v1LoginGroup)

	app.RegisterRoutes(router,
		homeGroup,
		v1LoginGroup,
		v1FreeGroup,
	)
}
