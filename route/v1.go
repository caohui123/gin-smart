package route

import (
	"github.com/jangozw/gin-smart/pkg/app"
)

// v1 不需要登陆
func registerV1Free(g *app.RouteGroup) {
	g.GET("/config", func(c *app.Context) app.Err {
		return nil
	})
}

// v1 需要登陆
func registerV1Login(g *app.RouteGroup) {
}
