package erro

const CodeDefault = 1

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
		code: CodeDefault,
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

func Fail(code int, msg ...string) E {
	var eMsg string
	if len(msg) > 0 {
		eMsg = msg[0]
	}
	return &errInfo{code: code, msg: eMsg}
}
