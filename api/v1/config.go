package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/pkg/app"
)

// 全局配置
func Config(c *gin.Context){
	ctx := app.Ctx(c)
	data := map[string]interface{}{
		"title":  "global config",
		"detail": "xxx",
	}
	ctx.Success(data)
	return
}
