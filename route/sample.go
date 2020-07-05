package route

import "github.com/jangozw/gin-smart/pkg/app"

func registerSampleFree(g app.RouteGroup) {
	g.GET("/sample", func(c *app.Context) app.Err {
		return nil
	})
}
