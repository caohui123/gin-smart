package app

import (
	"errors"
	"fmt"
)

// 错误接口
type Err interface {
	Error() string
	Code() int
	ErrMsg() string
}

// 错误信息
type errInfo struct {
	code   int
	errMsg string
}

func (e *errInfo) Msg(msg string, args ...interface{}) *errInfo {
	e.errMsg = fmt.Sprintf(msg, args...)
	return e
}

func (e *errInfo) Error() string {
	return e.errMsg
}

func (e *errInfo) Code() int {
	return e.code
}

func (e *errInfo) ErrMsg() string {
	return e.errMsg
}

func (e *errInfo) ToError() error {
	return errors.New(e.Error())
}

func Error(code int, msg interface{}, args ...interface{}) *errInfo {
	var errMsg string
	if v, ok := msg.(error); ok {
		errMsg = v.Error()
	} else {
		errMsg = fmt.Sprintf("%s", errMsg)
	}
	return &errInfo{code: code, errMsg: fmt.Sprintf(errMsg, args...)}
}

// 创建一个err
func CodeErr(code int) *errInfo {
	return &errInfo{code: code}
}

//  创建一个err
func ErrCode(code int) *errInfo {
	return &errInfo{code: code}
}
