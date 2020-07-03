package route

import (
	"gitlab.gosccap.cn/bourse/avian/pkg/app"
)

func registerRoot(g *app.RouteGroup) {
	g.GET("/test").HandlerFunc(func(c *app.Context) app.Err {
		c.Output("welcome!")
		return nil
	})
}
