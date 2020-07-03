package route

import "github.com/jangozw/gin-smart/pkg/app"

func registerSampleFree(g app.RouteGroup) {
	g.GET("/sample").HandlerFunc(func(c *app.Context) app.Err {
		c.Output("This is sample")
		return nil
	})
}
