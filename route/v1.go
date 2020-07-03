package route

import (
	v1 "gitlab.gosccap.cn/bourse/avian/api/v1"
	"gitlab.gosccap.cn/bourse/avian/pkg/app"
)

// v1 不需要登陆
func registerV1Free(g *app.RouteGroup) {
	g.GET("/config").HandlerFunc(v1.Config)
}

// v1 需要登陆
func registerV1Login(g *app.RouteGroup) {
}
