package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/erro"
	"net/http"

	"github.com/jangozw/gin-smart/pkg/app"
)

func CommonMiddleware(c *gin.Context) {
	if c.Request.Method != "GET" && c.GetHeader("Content-Type") != "application/json" {
		c.AbortWithStatusJSON(http.StatusOK, app.ResponseFailByCode(erro.InvalidContentTypeJSON))
		return
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Build-Info", app.BuildInfo)
	c.Next()
}
