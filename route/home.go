package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/pkg/app"
)

func registerRoot(r *gin.RouterGroup) {
	r.GET(`/home`, func(context *gin.Context) {
		ctx := app.Ctx(context)
		ctx.Success("Welcome")
	})
}
