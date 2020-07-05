package route

import (
	"github.com/jangozw/gin-smart/pkg/app"
)

func registerRoot(g *app.RouteGroup) {
	g.GET("/test", func(c *app.Context) app.Err {
		return nil
	})
}
