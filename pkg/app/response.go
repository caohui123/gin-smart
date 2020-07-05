package app

import (
	"time"
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
func ResponseFail(err Err) *response {
	return Response(err, struct{}{})
}

func ResponseFailByCode(code int) *response {
	return Response(Error(code), struct{}{})
}

func Response(err Err, data interface{}) *response {
	if err == nil {
		err = Error(Setting.ErrCodeSuccess)
	}
	errMsg := err.Msg()
	if errMsg == "" {
		errMsg = getCodeMsg(err.Code())
	}
	return &response{
		err.Code(),
		errMsg,
		time.Now().Unix(),
		data,
	}
}
