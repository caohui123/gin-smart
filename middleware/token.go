package middleware

import (
	"github.com/jangozw/gin-smart/erro"
	"github.com/jangozw/gin-smart/param"

	"github.com/gin-gonic/gin"

	"github.com/jangozw/gin-smart/pkg/app"
)

func CheckToken(c *gin.Context) {
	token := c.GetHeader(param.TokenHeaderKey)
	if token == "" {
		app.AbortJSON(app.Ctx(c), app.ResponseFailByCode(erro.ErrToken))
		return
	}
	if _, err := app.ParseUserByToken(token); err != nil {
		app.AbortJSON(app.Ctx(c), app.ResponseFailByCode(erro.ErrToken))
		return
	}
	// 继续下一步
	c.Next()
}
