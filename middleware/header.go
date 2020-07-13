package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/erro"

	"github.com/jangozw/gin-smart/pkg/app"
)

func Header(c *gin.Context) {
	if c.Request.Method != "GET" && c.GetHeader("Content-Type") != "application/json" {
		app.AbortJSON(app.Ctx(c), app.ResponseFailByCode(erro.InvalidContentTypeJSON))
		return
	}
	// 根据情况修改跨域
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Build-Info", app.BuildInfo)
	c.Next()
}
