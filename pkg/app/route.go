package app

import (
	"github.com/gin-gonic/gin"
)

type Route struct {
	Group  string
	Method string
	Path   string
	// Handler     ApiInstanceFunc
	HandlerFunc ApiHandlerFunc
}

func (r *Route) FullPath() string {
	return r.Group + r.Path
}

type RouteGroup struct {
	Name       string
	Middleware []gin.HandlerFunc
	Routes     []Route
}

// 通过对象实例处理请求
// type ApiInstanceFunc func() ApiIF

// 直接func 处理请求
type ApiHandlerFunc func(c *Context) Err

/*func (g *RouteGroup) Handler(instance ApiHandlerIF) {
	g.Routes[len(g.Routes)-1].HandlerInstance = instance
}*/

/*func (g *RouteGroup) HandlerFunc(fu ApiHandlerFunc) {
	g.Routes[len(g.Routes)-1].HandlerFunc = fu
}

func (g *RouteGroup) Handler(fu ApiInstanceFunc) {
	g.Routes[len(g.Routes)-1].Handler = fu
}*/

func (g *RouteGroup) Use(mid ...gin.HandlerFunc) *RouteGroup {
	g.Middleware = append(g.Middleware, mid...)
	return g
}

func (g *RouteGroup) GET(path string, handler ApiHandlerFunc) *RouteGroup {
	item := Route{Group: g.Name, Method: "GET", Path: path, HandlerFunc: handler}
	g.Routes = append(g.Routes, item)
	return g
}

func (g *RouteGroup) POST(path string, handler ApiHandlerFunc) *RouteGroup {
	item := Route{Group: g.Name, Method: "POST", Path: path, HandlerFunc: handler}
	g.Routes = append(g.Routes, item)
	return g
}

func (g *RouteGroup) PUT(path string, handler ApiHandlerFunc) *RouteGroup {
	item := Route{Group: g.Name, Method: "PUT", Path: path, HandlerFunc: handler}
	g.Routes = append(g.Routes, item)
	return g
}

func (g *RouteGroup) Any(path string, handler ApiHandlerFunc) *RouteGroup {
	item := Route{Group: g.Name, Method: "Any", Path: path, HandlerFunc: handler}
	g.Routes = append(g.Routes, item)
	return g
}

func (g *RouteGroup) DELETE(path string, handler ApiHandlerFunc) *RouteGroup {
	item := Route{Group: g.Name, Method: "DELETE", Path: path, HandlerFunc: handler}
	g.Routes = append(g.Routes, item)
	return g
}

func (g *RouteGroup) PATCH(path string, handler ApiHandlerFunc) *RouteGroup {
	item := Route{Group: g.Name, Method: "PATCH", Path: path, HandlerFunc: handler}
	g.Routes = append(g.Routes, item)
	return g
}

func (g *RouteGroup) OPTIONS(path string, handler ApiHandlerFunc) *RouteGroup {
	item := Route{Group: g.Name, Method: "OPTIONS", Path: path, HandlerFunc: handler}
	g.Routes = append(g.Routes, item)
	return g
}

func NewRouteGroup(name string) *RouteGroup {
	return &RouteGroup{
		Name:       name,
		Middleware: make([]gin.HandlerFunc, 0),
		Routes:     make([]Route, 0),
	}
}

func RegisterRoutes(engine *gin.Engine, groups ...*RouteGroup) {
	for _, g := range groups {
		gRoute := engine.Group(g.Name).Use(g.Middleware...)
		for _, r := range g.Routes {
			if err := routeManager(r); err != nil {
				continue
			}
			if r.Method == "GET" {
				gRoute = gRoute.GET(r.Path, APIHandler(r))
			} else if r.Method == "POST" {
				gRoute = gRoute.POST(r.Path, APIHandler(r))
			} else if r.Method == "Any" {
				gRoute = gRoute.Any(r.Path, APIHandler(r))
			} else if r.Method == "DELETE" {
				gRoute = gRoute.DELETE(r.Path, APIHandler(r))
			} else if r.Method == "PATCH" {
				gRoute = gRoute.PATCH(r.Path, APIHandler(r))
			} else if r.Method == "OPTIONS" {
				gRoute = gRoute.OPTIONS(r.Path, APIHandler(r))
			} else if r.Method == "PUT" {
				gRoute = gRoute.PUT(r.Path, APIHandler(r))
			} else {
				panic(r.Method + ": unknown route method")
			}
		}
	}
}
