package app

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/erro"
	"github.com/sirupsen/logrus"
)

var BuildInfo string // 编译的app版本

const (
	TimeFormatFullDate = "2006-01-02 15:04:05" // 常规类型
	EnvLocal           = "local"
	EnvDev             = "dev"
	EnvTest            = "test"
	EnvProduction      = "production"
	CtxStartTime       = "ctx-start-time"
)

type paramsCheck interface {
	Check() error
}

type TriggerIF interface {
	Do()
}

// 触发
func Trigger(tg TriggerIF) {
	go tg.Do()
}

// 环境
func IsEnvLocal() bool {
	return CurrentEnv() == EnvLocal
}

func IsEnvDev() bool {
	return CurrentEnv() == EnvDev
}

func IsEnvTest() bool {
	return CurrentEnv() == EnvTest
}

func IsEnvProduction() bool {
	return CurrentEnv() == EnvProduction
}

func CurrentEnv() string {
	if Cfg != nil {
		return Cfg.General.Env
	}
	return ""
}

// 对handler 执行捕获异常, 每一个中间件或api handler 都应该被warp
func Warp(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(CtxStartTime, time.Now())
		defer func() {
			if msg := recover(); msg != nil {
				LogApiPanic(c, msg)
				Ctx(c).Fail(erro.Inner(fmt.Sprintf("%v", msg)))
			}
		}()
		handler(c)
	}
}

// api 接口日志记录请求和返回
func LogApi(c *gin.Context) {
	ctx := Ctx(c)
	user, _ := ctx.LoginUser()

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
		query = queryPostToMap(postData)
	}

	// log 里有json.Marshal() 导致url转义字符
	Logger.WithFields(logrus.Fields{
		"uid":      user.ID,
		"query":    query,
		"response": resp,
	}).Infof("%s | %s |t=%3v | %s", c.Request.Method, c.Request.URL.RequestURI(), latency, c.ClientIP())
}

func LogApiPanic(c *gin.Context, panicMsg interface{}) {
	ctx := Ctx(c)
	user, _ := ctx.LoginUser()
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
		query = queryPostToMap(postData)
	}

	// log 里有json.Marshal() 导致url转义字符
	Logger.WithFields(logrus.Fields{
		"uid":      user.ID,
		"query":    query,
		"response": resp,
		"method":   c.Request.Method,
		"uri":      c.Request.URL.RequestURI(),
		"latency":  latency,
		"ip":       c.ClientIP(),
	}).Infof("--panic: %s | %s %s", panicMsg, c.Request.Method, c.Request.URL.RequestURI())
}

func queryPostToMap(data []byte) map[string]interface{} {
	m := make(map[string]interface{})
	json.Unmarshal(data, &m)
	return m
}
