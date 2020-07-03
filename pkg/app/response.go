package app

import (
	"fmt"
	"time"
)

// api 接口响应结果
type Response struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// 返回成功， 有数据
func ResponseSuccess(data interface{}) *Response {
	return ResponseWithCode(nil, data)
}

// 返回成功，不带数据，通用型
func ResponseSuccessSimple() *Response {
	return ResponseWithCode(nil, struct{}{})
}

// 返回错误， 带错误信息
func ResponseFail(err Err) *Response {
	return ResponseWithCode(err, struct{}{})
}

// 返回错误， 通用型
func ResponseFailByMsg(msg string, args ...interface{}) *Response {
	err := Error(errCodeFail, fmt.Sprintf(msg, args...))
	return ResponseWithCode(err, struct{}{})
}

func ResponseFailByCode(code int) *Response {
	return ResponseWithCode(CodeErr(code), struct{}{})
}

func ResponseWithCode(err Err, data interface{}) *Response {
	// 没有错误就是成功返回
	if err == nil {
		err = CodeErr(errCodeSuccess)
	}
	errMsg := err.ErrMsg()
	codMsg := getCodeMsg(err.Code())
	if errMsg == "" {
		errMsg = codMsg
	} else if err.Code() != errCodeFail {
		errMsg = codMsg + ":" + errMsg
	}
	return &Response{
		err.Code(),
		errMsg,
		time.Now().Unix(),
		data,
	}
}
