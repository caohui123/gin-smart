package middleware

import (
	"net/http"

	"github.com/jangozw/gin-smart/erro"
	"github.com/jangozw/gin-smart/param"

	"github.com/gin-gonic/gin"

	"github.com/jangozw/gin-smart/pkg/app"
)

func CheckToken(c *gin.Context) {
	abort := func() {
		c.AbortWithStatusJSON(http.StatusOK, app.ResponseFailByCode(erro.ErrToken))
		return
	}
	token := c.GetHeader(param.TokenHeaderKey)
	if token == "" {
		abort()
		return
	}
	if _, err := app.ParseUserByToken(token); err != nil {
		abort()
		return
	}
	// 继续下一步
	c.Next()
}
