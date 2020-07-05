package app

import (
	"fmt"
	"time"
)

const (
	errCodeSuccess             = 200
	errCodeFail                = 400
	errCodeInvalidLoginUser    = 5501
	errCodeInvalidRequestParam = 5502
	errCodeNotSetRouteHandler  = 5503
)

var innerErrCodeMap = map[int]string{
	errCodeSuccess:             "请求成功",
	errCodeFail:                "请求失败",
	errCodeInvalidLoginUser:    "无效的登陆用户",
	errCodeInvalidRequestParam: "请求参数验证失败",
	errCodeNotSetRouteHandler:  "路由没有绑定实例",
}

type SettingFields struct {
	ErrCodeMap   map[int]string
	Boot         string
	BuildAt      string
	StartAt      time.Time
	BuildVersion string
	StartArgs    StarArgs
}

var Setting SettingFields

func getCodeMsg(code int) string {
	if v, ok := Setting.ErrCodeMap[code]; ok {
		return v
	}
	if v, ok := innerErrCodeMap[code]; ok {
		return v
	}
	return fmt.Sprintf("facility:unknown code:%d", code)
}
