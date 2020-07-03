package route

import (
	v1 "github.com/jangozw/gin-smart/api/v1"
	"github.com/jangozw/gin-smart/pkg/app"
)

// v1 不需要登陆
func registerV1Free(g *app.RouteGroup) {
	g.GET("/config").HandlerFunc(v1.Config)
}

// v1 需要登陆
func registerV1Login(g *app.RouteGroup) {
}
