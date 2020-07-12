package app

import (
	"time"

	"github.com/jangozw/gin-smart/erro"
)

// api 接口响应结果
type response struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// 返回成功， 有数据
func ResponseSuccess(data interface{}) *response {
	return Response(nil, data)
}

// 返回错误， 带错误信息
func ResponseFail(err erro.E) *response {
	return Response(err, struct{}{})
}

func ResponseFailByCode(code int) *response {
	return Response(erro.Info(code), struct{}{})
}

func Response(err erro.E, data interface{}) *response {
	if err == nil {
		err = erro.Info(erro.Success)
	}
	errMsg := err.Msg()
	if errMsg == "" {
		errMsg = erro.GetCodeMsg(err.Code())
	}
	return &response{
		err.Code(),
		errMsg,
		time.Now().Unix(),
		data,
	}
}
