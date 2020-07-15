package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/erron"
)

// api 处理函数类型
type ApiHandlerFunc func(c *gin.Context) (data interface{}, err erron.E)

// 注册路由函数类型
type RegisterRouteFunc func(engine *Engine)

// 定义路由 gin.Engine 统一加warp
func NewGin(registerRoutes RegisterRouteFunc) *Engine {
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
		defer func() {
			if msg := recover(); msg != nil {
				defer func() {
					if msg := recover(); msg != nil {
						c.AbortWithStatusJSON(http.StatusOK, response{Code: erron.ErrInternal, Msg: "bad request", Timestamp: time.Now().Unix(), Data: nil})
					}
				}()
				err := erron.Inner(fmt.Sprintf("%v", msg))
				AbortJSON(c, ResponseFail(err))
				LogApiPanic(c, msg)
			}
		}()
		data, err := handler(c)
		OutputJSON(c, Response(err, data))
	}
}

// 中间件 捕获异常
func WarpMiddleware(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if msg := recover(); msg != nil {
				defer func() {
					if msg := recover(); msg != nil {
						c.AbortWithStatusJSON(http.StatusOK, response{Code: erron.ErrInternal, Msg: "bad request", Timestamp: time.Now().Unix(), Data: nil})
					}
				}()
				err := erron.Inner(fmt.Sprintf("%v", msg))
				AbortJSON(c, ResponseFail(err))
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
func OutputJSON(c *gin.Context, resp *response) {
	setResponse(c, resp)
	c.JSON(http.StatusOK, resp)
}

// 中断并输出
func AbortJSON(c *gin.Context, resp *response) {
	setResponse(c, resp)
	c.AbortWithStatusJSON(http.StatusOK, resp)
}

// api 请求发生了panic 记入日志
func LogApiPanic(c *gin.Context, panicMsg interface{}) {
	user, _ := GetLoginUser(c)
	start := c.GetTime(CtxStartTime)
	// 执行时间
	latency := time.Now().Sub(start)
	resp, ok := c.Get(CtxKeyResponse)
	if !ok {
		resp = struct{}{}
	}
	var query interface{}
	if c.Request.Method == "GET" {
		query = c.Request.URL.Query()
	} else {
		postData, _ := c.GetRawData()
		query := make(map[string]interface{})
		json.Unmarshal(postData, &query)
	}

	// log 里有json.Marshal() 导致url转义字符
	Logger.WithFields(logrus.Fields{
		"uid":      user.ID,
		"query":    query,
		"response": resp,
		"method":   c.Request.Method,
		"uri":      c.Request.URL.RequestURI(),
		"latency":  fmt.Sprintf("%3v", latency),
		"ip":       c.ClientIP(),
	}).Infof("--panic: %s | %s %s", panicMsg, c.Request.Method, c.Request.URL.RequestURI())
}
