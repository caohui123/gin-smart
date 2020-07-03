package app

type ApiIF interface {
	Prepare(ctx *Context) Err
	Handler() Err
}
