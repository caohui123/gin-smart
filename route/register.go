package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/middleware"
)

func Register(router *gin.Engine) {
	// 用户token验证接口实例

	sampleGroup := router.Group("/sample")
	homeGroup := router.Group("/")
	v1LoginGroup := router.Group("/v1", middleware.LogToFile(), middleware.CommonMiddleware, middleware.CheckToken())
	v1FreeGroup := router.Group("/v1", middleware.LogToFile(), middleware.CommonMiddleware)
	registerRoot(homeGroup)
	registerSampleFree(sampleGroup)
	registerV1Free(v1FreeGroup)
	registerV1Login(v1LoginGroup)
}
