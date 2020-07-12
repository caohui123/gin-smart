package app

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/erro"
	"github.com/jangozw/gin-smart/param"
	"github.com/jangozw/gin-smart/pkg/auth"
	"github.com/jangozw/gin-smart/pkg/util"
)

const (
	CtxKeyResponse = "ctx-response"
)

var inputTypeErr = errors.New("input 必须是一个结构体变量的地址")

// 登陆用户信息, 可根据业务需要扩充字段
type LoginUser struct {
	ID int64
}

// 根据token解析登陆用户
func ParseUserByToken(token string) (LoginUser, error) {
	user := LoginUser{}
	if token == "" {
		return user, errors.New("empty token")
	}
	if jwtPayload, err := auth.ParseJwtToken(token, Cfg.General.JwtSecret); err != nil {
		return user, err
	} else if err := jwtPayload.ParseUser(&user); err != nil {
		return user, err
	}
	if user.ID == 0 {
		return user, errors.New("invalid login user")
	}
	return user, nil
}

// 用户api 请求
type Context struct {
	*gin.Context
	pager        *util.Pager
	loginUser    LoginUser
	outputPager  bool
	recordsCount uint
}

func Ctx(c *gin.Context) *Context {
	return &Context{Context: c}
}

// 展示分页
func (c *Context) Pager() *util.Pager {
	if c.pager == nil {
		pager := &util.Pager{}
		if err := c.ShouldBind(pager); err == nil {
			pager = util.NewRequestPager(pager.Page, pager.PageSize)
		}
		c.pager = pager
	}
	return c.pager
}

func (c *Context) LoginUser() (LoginUser, error) {
	if c.loginUser.ID == 0 {
		user, err := ParseUserByToken(c.GetHeader(param.TokenHeaderKey))
		if err != nil {
			return c.loginUser, err
		}
		c.loginUser = user
	}
	return c.loginUser, nil
}

func (c *Context) MustLoginUser() LoginUser {
	user, err := c.LoginUser()
	if err != nil || user.ID == 0 {
		panic(err)
	}
	return c.loginUser
}

// 绑定输入参数
func (c *Context) BindInput(input interface{}) error {
	if err := checkInput(input); err == nil {
		if err := c.ShouldBind(input); err != nil {
			return err
		}
		// 如果实现了 params 接口就验证参数
		if obj, ok := input.(paramsCheck); ok {
			if err := obj.Check(); err != nil {
				return err
			}
		}
	}
	return nil
}

// 返回成功， 有数据
func (c *Context) Success(data interface{}) {
	c.Response(nil, data)
}

// 返回错误， 通用型
func (c *Context) Fail(err erro.E) {
	if err == nil {
		// err = erro.Info(-1)
	}
	c.Response(err, struct{}{})
}

func (c *Context) Response(err erro.E, data interface{}) {
	output := Response(err, c.adjustOutput(data))
	c.Set(CtxKeyResponse, output)
	c.JSON(http.StatusOK, output)
}

func (c *Context) adjustOutput(data interface{}) interface{} {
	// 如果是分页请求,改变下返回结构
	if c.outputPager == true {
		type pageOutput struct {
			Pager util.PageResp `json:"pager"`
			List  interface{}   `json:"list"`
		}
		data = pageOutput{
			Pager: util.PageResp{
				Page:     c.Pager().Page,
				PageSize: c.Pager().PageSize,
				Total:    c.recordsCount,
			},
			List: data,
		}
	}
	return data
}

// 输出结构展示分页
func (c *Context) SetPager(count uint) {
	c.recordsCount = count
	c.outputPager = true
}

func checkInput(input interface{}) error {
	rv := reflect.ValueOf(input)
	if rv.Kind() != reflect.Ptr || rv.IsNil() || !rv.IsValid() {
		return inputTypeErr
	}
	rt := reflect.TypeOf(input)
	if rt.Kind() == reflect.Ptr {
		if rt.Elem().Kind() != reflect.Struct {
			return inputTypeErr
		}
	} else {
		if rt.Kind() != reflect.Struct {
			return inputTypeErr
		}
	}
	return nil
}
