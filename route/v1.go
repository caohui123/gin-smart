package route

import (
	"github.com/gin-gonic/gin"
)

// v1 不需要登陆
func registerV1Free(r *gin.RouterGroup) {
	r.GET(`/config`, func(context *gin.Context) {
	})
}

// v1 需要登陆
func registerV1Login(r *gin.RouterGroup) {
}
