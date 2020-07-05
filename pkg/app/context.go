package app

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

const (
	CtxKeyLoginUser = "ctx-login-user"
	CtxKeyNeedLogin = "ctx-need-login"
	CtxKeyResponse  = "ctx-response"
)

var (
	inputTypeErr = errors.New("input 必须是一个结构体变量的地址")
)

// 登陆用户接口
type LoginIF interface {
	User() LoginUser
}

// 登陆用户信息, 可根据业务需要扩充字段
type LoginUser struct {
	ID    uint
	Extra interface{}
}

func (l *LoginUser) User() LoginUser {
	return *l
}

// info 要传地址
func (l *LoginUser) ExtraTo(info interface{}) error {
	bytes, err := json.Marshal(l.Extra)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, info)
}

// 用户api 请求
type Context struct {
	*gin.Context
	pager        *Pager
	loginUser    *LoginUser
	outputPager  bool
	recordsCount uint
}

func Ctx(c *gin.Context) *Context {
	return &Context{Context: c}
}

// 展示分页
func (c *Context) Pager() *Pager {
	if c.pager == nil {
		pager := &Pager{}
		if err := c.ShouldBind(pager); err == nil {
			pager = NewRequestPager(pager.Page, pager.PageSize)
		}
		c.pager = pager
	}
	return c.pager
}

func (c *Context) LoginUser() *LoginUser {
	if c.loginUser == nil {
		loginUser := LoginUser{}
		if c.GetBool(CtxKeyNeedLogin) == true {
			if obj, exists := c.Get(CtxKeyLoginUser); exists {
				loginUser = obj.(LoginUser)
			}
		}
		c.loginUser = &loginUser
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
func (c *Context) Fail(err Err) {
	if err == nil {
		err = Error(-1)
	}
	c.Response(err, struct{}{})
}

func (c *Context) Response(err Err, data interface{}) {
	output := Response(err, c.adjustOutput(data))
	c.Set(CtxKeyResponse, output)
	c.JSON(http.StatusOK, output)
}

func (c *Context) adjustOutput(data interface{}) interface{} {
	// 如果是分页请求,改变下返回结构
	if c.outputPager == true {
		type pageOutput struct {
			Pager pagerResp   `json:"pager"`
			List  interface{} `json:"list"`
		}
		data = pageOutput{
			Pager: pagerResp{
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
