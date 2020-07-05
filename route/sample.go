package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/pkg/app"
)

func registerSampleFree(r *gin.RouterGroup) {
	r.GET("", func(context *gin.Context) {
		ctx := app.Ctx(context)
		ctx.Success("Sample here")
	})

	r.GET("/list", func(context *gin.Context) {
		ctx := app.Ctx(context)

		// 请求的分页信息
		/*
			page := ctx.Pager().Page
			pageSize := ctx.Pager().PageSize
			offset := ctx.Pager().Offset()
			limit := ctx.Pager().Limit()
		*/

		type user struct {
			Name string `json:"name"`
			Age  int8   `json:"age"`
		}
		list := make([]user, 0)
		list = append(list, user{
			Name: "Dog",
			Age:  2,
		})
		list = append(list, user{
			Name: "Cat",
			Age:  3,
		})

		var total uint = 2
		// 设置数据总数， 自动调整输出结构为分页结构
		ctx.SetPager(total)
		ctx.Success(list)
	})
}
