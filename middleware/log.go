package middleware

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/pkg/app"
	"github.com/jangozw/go-api-facility/auth"
	"github.com/sirupsen/logrus"
)

func LogToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// 处理请求
		c.Next()
		// 写日志的时间单独一个goroutine 处理，不占用接口调用时间
		go func() {
			// 执行时间
			latency := time.Now().Sub(start)
			resp, ok := c.Get(app.CtxKeyResponse)
			if !ok {
				resp = struct{}{}
			}

			var uid string
			token := c.GetHeader(TokenHeaderKey)
			if token != "" {
				if payload, err := auth.ParseJwtToken(token, app.Runner.Cfg.Encrypt.JwtSecret); err == nil {
					uid = payload.AccountID
				}
			}
			var query interface{}
			if c.Request.Method == "GET" {
				query = c.Request.URL.Query()
			} else {
				postData, _ := c.GetRawData()
				query = queryPostToMap(postData)
			}

			// uri := strings.ReplaceAll(c.Request.URL.RequestURI(), "\\u0026", "&")
			uri := c.Request.URL.RequestURI()
			// log 里有json.Marshal() 导致url转义字符
			app.Runner.Log.Api.WithFields(logrus.Fields{
				"uid":      uid,
				"query":    query,
				"response": resp,
			}).Infof("%s | %s |t=%3v | %s", c.Request.Method, uri, latency, c.ClientIP())
		}()
	}
}

func queryPostToMap(data []byte) map[string]interface{} {
	m := make(map[string]interface{})
	json.Unmarshal(data, &m)
	return m
}
