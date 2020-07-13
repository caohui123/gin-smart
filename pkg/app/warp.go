package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/erro"
	"github.com/jangozw/gin-smart/pkg/util"
	"net/http"
	"time"
)

type ApiHandlerFunc func(c *Context) (data interface{}, err erro.E)

// 定义路由 gin.Engine 统一加warp
func Engine(engine *gin.Engine) *EngineWarp {
	return &EngineWarp{engine: engine}
}

type RouteGroup struct {
	rg *gin.RouterGroup
}

// 需要什么方法自由搬运 gin.RouteGroup
func (r *RouteGroup) GET(relativePath string, handler ApiHandlerFunc) {
	r.rg.GET(relativePath, WarpApi(handler))
}
func (r *RouteGroup) POST(relativePath string, handler ApiHandlerFunc) {
	r.rg.POST(relativePath, WarpApi(handler))
}
func (r *RouteGroup) DELETE(relativePath string, handler ApiHandlerFunc) {
	r.rg.DELETE(relativePath, WarpApi(handler))
}
func (r *RouteGroup) Any(relativePath string, handler ApiHandlerFunc) {
	r.rg.Any(relativePath, WarpApi(handler))
}
func (r *RouteGroup) PATCH(relativePath string, handler ApiHandlerFunc) {
	r.rg.PATCH(relativePath, WarpApi(handler))
}
func (r *RouteGroup) OPTIONS(relativePath string, handler ApiHandlerFunc) {
	r.rg.OPTIONS(relativePath, WarpApi(handler))
}
func (r *RouteGroup) PUT(relativePath string, handler ApiHandlerFunc) {
	r.rg.PUT(relativePath, WarpApi(handler))
}
func (r *RouteGroup) HEAD(relativePath string, handler ApiHandlerFunc) {
	r.rg.HEAD(relativePath, WarpApi(handler))
}

// 需要用其他的再加

//
type EngineWarp struct {
	engine *gin.Engine
}

func (e *EngineWarp) Group(relativePath string, handlers ...gin.HandlerFunc) *RouteGroup {
	for i, h := range handlers {
		handlers[i] = WarpMiddleware(h)
	}
	group := e.engine.Group(relativePath, handlers...)
	return &RouteGroup{rg: group}
}

// 需要用其他的再加

// api 捕获异常
func WarpApi(handler ApiHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{Context:c}
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
				AbortJSON(&Context{Context:c}, ResponseFail(err))
				LogApiPanic(c, msg)
			}
		}()
		if _, ok := c.Get(CtxStartTime); !ok {
			c.Set(CtxStartTime, time.Now())
		}
		handler(c)
	}
}

// 输出
func OutputJSON(c *Context, resp *response) {
	// 如果是分页请求,改变下返回结构
	if c.outputPager == true {
		type pageOutput struct {
			Pager util.PageResp `json:"pager"`
			List  interface{}   `json:"list"`
		}
		resp.Data = pageOutput{
			Pager: util.PageResp{
				Page:     c.Pager().Page,
				PageSize: c.Pager().PageSize,
				Total:    c.recordsCount,
			},
			List: resp,
		}
	}
	c.Set(CtxKeyResponse, resp)
	c.JSON(http.StatusOK, resp)
}

func AbortJSON(c *Context, resp *response) {
	c.Set(CtxKeyResponse, resp)
	c.AbortWithStatusJSON(http.StatusOK, resp)
}
