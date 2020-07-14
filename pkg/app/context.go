package app

import (
	"errors"
	"reflect"

	"github.com/jangozw/gin-smart/erro"

	"github.com/gin-gonic/gin"
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

type TokenPayload struct {
	UserID int64 `json:"uid"`
}

// 根据token解析登陆用户
func ParseUserByToken(token string) (TokenPayload, error) {
	user := TokenPayload{}
	if token == "" {
		return user, errors.New("empty token")
	}
	if jwtPayload, err := auth.ParseJwtToken(token, Cfg.General.JwtSecret); err != nil {
		return user, err
	} else if err := jwtPayload.ParseUser(&user); err != nil {
		return user, err
	}
	if user.UserID == 0 {
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

func (c *Context) LoginUser() (LoginUser, error) {
	if c.loginUser.ID == 0 {
		info, err := ParseUserByToken(c.GetHeader(param.TokenHeaderKey))
		if err != nil {
			return c.loginUser, err
		}
		c.loginUser = LoginUser{ID: info.UserID}
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
func (c *Context) BindInput(input interface{}) erro.E {
	if err := checkInput(input); err == nil {
		if err := c.ShouldBind(input); err != nil {
			return erro.Fail(erro.ErrRequestParams, err.Error())
		}
		// 如果实现了 params 接口就验证参数
		if obj, ok := input.(paramsCheck); ok {
			if err := obj.Check(); err != nil {
				return erro.Fail(erro.ErrRequestParams, err.Error())
			}
		}
	}
	return nil
}

// 输出结构展示分页
func (c *Context) SetPager(count uint) {
	c.recordsCount = count
	c.outputPager = true
}

// 展示的分页
func (c *Context) GetPager() *Pager {
	pager := &Pager{}
	c.ShouldBind(pager)
	pager.Secure()
	return pager
}

func (c *Context) setResponse(resp *response) {
	c.Set(CtxKeyResponse, resp)
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
