package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/erro"
)

// api 处理函数类型
type ApiHandlerFunc func(c *Context) (data interface{}, err erro.E)

// 注册路由函数类型
type RegisterRouteFunc func(engine *Engine)

// 定义路由 gin.Engine 统一加warp
func NewEngine(registerRoutes RegisterRouteFunc) *Engine {
	// 注册路由
	if IsEnvLocal() || IsEnvDev() {
		gin.SetMode(gin.DebugMode)
	}
	eng := &Engine{gin.New()}
	registerRoutes(eng)
	return eng
}

type routeGroup struct {
	rg *gin.RouterGroup
}

// 路由组中间件
func (r *routeGroup) Use(middleware ...gin.HandlerFunc) *routeGroup {
	for i, h := range middleware {
		middleware[i] = WarpMiddleware(h)
	}
	r.rg.Use(middleware...)
	return r
}

// 需要什么方法自由搬运 gin.routeGroup
func (r *routeGroup) GET(relativePath string, handler ApiHandlerFunc) {
	r.rg.GET(relativePath, WarpApi(handler))
}

func (r *routeGroup) POST(relativePath string, handler ApiHandlerFunc) {
	r.rg.POST(relativePath, WarpApi(handler))
}

func (r *routeGroup) DELETE(relativePath string, handler ApiHandlerFunc) {
	r.rg.DELETE(relativePath, WarpApi(handler))
}

func (r *routeGroup) Any(relativePath string, handler ApiHandlerFunc) {
	r.rg.Any(relativePath, WarpApi(handler))
}

func (r *routeGroup) PATCH(relativePath string, handler ApiHandlerFunc) {
	r.rg.PATCH(relativePath, WarpApi(handler))
}

func (r *routeGroup) OPTIONS(relativePath string, handler ApiHandlerFunc) {
	r.rg.OPTIONS(relativePath, WarpApi(handler))
}

func (r *routeGroup) PUT(relativePath string, handler ApiHandlerFunc) {
	r.rg.PUT(relativePath, WarpApi(handler))
}

func (r *routeGroup) HEAD(relativePath string, handler ApiHandlerFunc) {
	r.rg.HEAD(relativePath, WarpApi(handler))
}

// 需要用其他的再加

//
type Engine struct {
	engine *gin.Engine
}

func (e *Engine) Group(relativePath string, handlers ...gin.HandlerFunc) *routeGroup {
	for i, h := range handlers {
		handlers[i] = WarpMiddleware(h)
	}
	group := e.engine.Group(relativePath, handlers...)
	return &routeGroup{rg: group}
}

func (e *Engine) Run() error {
	return e.engine.Run(HttpServeAddr())
}

// 需要用其他的再加

// api 捕获异常
func WarpApi(handler ApiHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{Context: c}
		defer func() {
			if msg := recover(); msg != nil {
				defer func() {
					if msg := recover(); msg != nil {
						c.AbortWithStatusJSON(http.StatusOK, response{Code: erro.ErrInternal, Msg: "bad request", Timestamp: time.Now().Unix(), Data: nil})
					}
				}()
				err := erro.Inner(fmt.Sprintf("%v", msg))
				AbortJSON(ctx, ResponseFail(err))
				LogApiPanic(c, msg)
			}
		}()
		data, err := handler(ctx)
		OutputJSON(ctx, Response(err, data))
	}
}

// 中间件 捕获异常
func WarpMiddleware(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if msg := recover(); msg != nil {
				defer func() {
					if msg := recover(); msg != nil {
						c.AbortWithStatusJSON(http.StatusOK, response{Code: erro.ErrInternal, Msg: "bad request", Timestamp: time.Now().Unix(), Data: nil})
					}
				}()
				err := erro.Inner(fmt.Sprintf("%v", msg))
				AbortJSON(&Context{Context: c}, ResponseFail(err))
				LogApiPanic(c, msg)
			}
		}()
		if _, ok := c.Get(CtxStartTime); !ok {
			c.Set(CtxStartTime, time.Now())
		}
		handler(c)
	}
}

// 正常输出
func OutputJSON(c *Context, resp *response) {
	c.setResponse(resp)
	c.JSON(http.StatusOK, resp)
}

// 中断并输出
func AbortJSON(c *Context, resp *response) {
	c.setResponse(resp)
	c.AbortWithStatusJSON(http.StatusOK, resp)
}
