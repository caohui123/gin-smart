package app

// 错误接口
type Err interface {
	Error() string
	Code() int
	Msg() string
}

// 错误信息
type errInfo struct {
	code int
	msg  string
}

func (e *errInfo) Error() string {
	return e.msg
}

func (e *errInfo) Code() int {
	return e.code
}

func (e *errInfo) Msg() string {
	return e.msg
}

func Error(code int, msg ...string) *errInfo {
	var eMsg string
	if len(msg) > 0 {
		eMsg = msg[0]
	}
	return &errInfo{code: code, msg: eMsg}
}
