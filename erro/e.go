package erro

// 错误接口
type E interface {
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

func New(text string) E {
	return &errInfo{
		code: Failed,
		msg:  text,
	}
}

func Info(code int, msg ...string) E {
	var eMsg string
	if len(msg) > 0 {
		eMsg = msg[0]
	}
	return &errInfo{code: code, msg: eMsg}
}

func Inner(msg ...string) E {
	var eMsg string
	if len(msg) > 0 {
		eMsg = msg[0]
	}
	return &errInfo{code: ErrInternal, msg: eMsg}
}

func Fail(code int, msg ...string) E {
	var eMsg string
	if len(msg) > 0 {
		eMsg = msg[0]
	}
	return &errInfo{code: code, msg: eMsg}
}

// FailCn Cn means Common
func FailCn() E {
	return &errInfo{code: Failed, msg: ""}
}

func FailBy(err error) E {
	return &errInfo{code: Failed, msg: err.Error()}
}
