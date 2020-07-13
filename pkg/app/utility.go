package app

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

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
