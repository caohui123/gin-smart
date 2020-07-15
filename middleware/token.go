package middleware

import (
	"github.com/jangozw/gin-smart/erron"
	"github.com/jangozw/gin-smart/param"

	"github.com/gin-gonic/gin"

	"github.com/jangozw/gin-smart/pkg/app"
)

func NeedLogin(c *gin.Context) {
	token := c.GetHeader(param.TokenHeaderKey)
	if token == "" {
		app.AbortJSON(c, app.ResponseFailByCode(erron.ErrToken))
		return
	}
	if _, err := app.ParseUserByToken(token); err != nil {
		app.AbortJSON(c, app.ResponseFail(erron.Info(erron.ErrToken, err.Error())))
		return
	}
	// 继续下一步
	c.Next()
}
