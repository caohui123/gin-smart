package app

import (
	"time"

	"github.com/jangozw/gin-smart/erron"
)

type responseWithPager struct {
	Pager *Pager      `json:"pager"`
	List  interface{} `json:"list"`
}

// 带分页的输出结果
func PagerResponse(pager *Pager, list interface{}) *responseWithPager {
	return &responseWithPager{
		Pager: pager,
		List:  list,
	}
}

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
func ResponseFail(err erron.E) *response {
	return Response(err, struct{}{})
}

func ResponseFailByCode(code int) *response {
	return Response(erron.Info(code), struct{}{})
}

func Response(err erron.E, data interface{}) *response {
	if err == nil {
		err = erron.Info(erron.Success)
	}
	errMsg := err.Msg()
	if errMsg == "" {
		errMsg = erron.GetCodeMsg(err.Code())
	}
	return &response{
		err.Code(),
		errMsg,
		time.Now().Unix(),
		data,
	}
}
