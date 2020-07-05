package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/errcode"

	"github.com/jangozw/gin-smart/pkg/app"
)

func CommonMiddleware(c *gin.Context) {
	if c.Request.Method != "GET" && c.GetHeader("Content-Type") != "application/json" {
		c.AbortWithStatusJSON(http.StatusOK, app.ResponseFailByCode(errcode.InvalidContentTypeJSON))
		return
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Start-Info", fmt.Sprintf("%s,%s,%s", app.Setting.BuildVersion, app.Setting.BuildAt, time.Now().Format(app.TimeFormatFullDate)))
	c.Next()
}
