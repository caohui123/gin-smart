package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/api/sample"
	"github.com/jangozw/gin-smart/pkg/app"
)

func registerSampleFree(r *gin.RouterGroup) {
	r.GET("", func(context *gin.Context) {
		ctx := app.Ctx(context)
		ctx.Success("Sample here")
	})

	r.GET("sample", sample.Sample)

	r.GET("/list", sample.List)
}
