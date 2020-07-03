# gin-smart

以gin为基础轻量封装的适合大多数项目通用的api开发套件




# 功能特色
* 使用gin框架, 文档： https://github.com/gin-gonic/gin
* 使用gorm数据操作, 文档： http://gorm.book.jasperxu.com 
* 使用JWT生成token, 结合redis双重验证。 文档 http://jwt.io
* 使用validator.v8验证器, 文档: https://godoc.org/gopkg.in/go-playground/validator.v8
* 使用logrus日志, 文档: https://github.com/sirupsen/logrus




# 使用
## 安装gofresh实时监测代码变动即时编译
参考：https://github.com/jangozw/gofresh

## 测试API接口

```bash
curl http://127.0.0.1/test
```

# 定义接口

1,  /route 里定义



2, /api 里写handler 

handler 有两种实现, 可以写成对象也可直接函数 示范：


```text
    // 定义一个路由
	g.GET("/config").HandlerFunc(v1.Config)

```

使用handler 函数:

```go
package v1

import "gitlab.gosccap.cn/bourse/avian/pkg/app"

func Config(c *app.Context) app.Err {
	var data = map[string]interface{}{
		"title":  "global config",
		"detail": "xxx",
	}
	c.Output(data)
	return nil
}

```


或者使用handler对象：
```go
package v1

import (
	"gitlab.gosccap.cn/bourse/avian/errcode"
	"gitlab.gosccap.cn/bourse/avian/param"
	"gitlab.gosccap.cn/bourse/avian/pkg/app"
)

// 路由
type SampleUserApiInstance struct {
	// 请求context
	ctx *app.Context

	// 请求参数自动解析
	input param.SampleUser

	// 返回正确的结果
	output param.SampleUserResponse
}

//
func (api *SampleUserApiInstance) Prepare(ctx *app.Context) app.Err {
	api.ctx = ctx
	if err := api.ctx.MustBinds(&api.input, &api.output); err != nil {
		return app.Error(errcode.InvalidBindValue, err.Error())
	}
	return nil
}

// 处理请求
func (api *SampleUserApiInstance) Handler() app.Err {
	// your code here...
	api.output.Name = "Sample2"
	api.output.Mobile = "16666666666666"
	return nil
}
```






