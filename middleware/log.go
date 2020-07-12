package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/pkg/app"
)

func LogRequest(c *gin.Context) {
	if _, ok := c.Get(app.CtxStartTime); !ok {
		c.Set(app.CtxStartTime, time.Now())
	}
	// 处理请求
	c.Next()
	// 写日志的时间单独一个goroutine 处理，不占用接口调用时间
	app.LogApi(c)
}
