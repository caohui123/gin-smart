# gin-smart

以gin为基础轻量封装的适合大多数项目通用的api开发套件




# 功能特色
* 使用gin框架, 文档： https://github.com/gin-gonic/gin
* 使用gorm数据操作, 文档： http://gorm.book.jasperxu.com 
* 使用JWT生成token, 结合redis双重验证。 文档 http://jwt.io
* 使用validator.v8验证器, 文档: https://godoc.org/gopkg.in/go-playground/validator.v8
* 使用logrus日志, 文档: https://github.com/sirupsen/logrus



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
	g.GET("/login").HandlerFunc(v1.SampleLogin)

```

使用handler 函数:

```go
package v1
func SampleLogin(c *app.Context) app.Err {
	var p param.LoginRequest
	if err := c.BindInput(&p); err != nil {
		return app.ErrCode(errcode.ErrRequestParams)
	}
	jwtToken, err := auth.Login(&service.JwtUserLogin{
		AccountID:  p.Mobile,
		AccountPwd: p.Pwd,
	})
	if err != nil {
		return app.ErrCode(errcode.Failed)
	}
	output := param.LoginResponse{Token: string(jwtToken)}
	c.Output(output)
	return nil
}

```





