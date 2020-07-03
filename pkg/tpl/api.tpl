package {package}

import (
	"{module}/param"
	"{module}/errcode"
	"{module}/pkg/app"
)

// {description}
type {structName} struct {
	// 请求context
	ctx *app.Context

	// 请求参数，自动解析到input
	input {input}

	// 请求返回结果赋值到output
	output {output}
}

// Prepare
// 请求的准备工作，传入context， 绑定输入输出参数
func (api *{structName}) Prepare(ctx *app.Context) app.Err {
	api.ctx = ctx
   	if err := api.ctx.MustBinds(&api.input, &api.output); err != nil {
   		return app.Error(errcode.InvalidBindValue, err.Error())
   	}
	return nil
}

// Handler
// 这里写你的业务逻辑...
// 只要将 api.output 赋值即可返回数据
func (api *{structName}) Handler() app.Err {
    // code here ...
	return nil
}

// 常用事项
// api.Pager 分页数据
// api.PagerOn(100) 设置数据总条数，即可返回带分页结构的数据
// app.Trigger(v) 设置一个触发器

