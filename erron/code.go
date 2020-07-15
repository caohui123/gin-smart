package erron

import "fmt"

// 1-1000 系统自带

const (
	Success                = 200
	Failed                 = 400
	ErrToken               = 401
	ErrInternal            = 500
	ErrRequestParams       = 600
	InvalidLoginUser       = 603
	InvalidContentTypeJSON = 604
	InvalidRouteHandler    = 605
	InvalidBindValue       = 606
)

// 自带
var appCodeMap = map[int]string{
	Success:                "请求成功",
	Failed:                 "请求失败",
	ErrToken:               "token错误",
	ErrInternal:            "内部错误或异常",
	ErrRequestParams:       "请求参数错误",
	InvalidLoginUser:       "获取登陆用户失败",
	InvalidContentTypeJSON: "Content-Type must be application/json",
	InvalidRouteHandler:    "路由未绑定实例",
	InvalidBindValue:       "绑定的入参或出参有误",
}

// 所有 code map 汇集计算到一起

var allCodeMaps = []map[int]string{
	appCodeMap,
}

func CodeMap() map[int]string {
	allMap := make(map[int]string)
	for _, item := range allCodeMaps {
		for k, v := range item {
			allMap[k] = v
		}
	}
	return allMap
}

func GetCodeMsg(code int) string {
	if v, ok := CodeMap()[code]; ok {
		return v
	}
	return fmt.Sprintf("unknown code:%d", code)
}
