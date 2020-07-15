package app

import (
	"errors"
	"reflect"

	"github.com/jangozw/gin-smart/erron"

	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/param"
	"github.com/jangozw/gin-smart/pkg/auth"
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

func GetLoginUser(c *gin.Context) (user LoginUser, err error) {
	info, err := ParseUserByToken(c.GetHeader(param.TokenHeaderKey))
	if err != nil {
		return user, err
	}
	return LoginUser{ID: info.UserID}, nil
}

func MustGetLoginUser(c *gin.Context) LoginUser {
	user, err := GetLoginUser(c)
	if err != nil || user.ID == 0 {
		panic(err)
	}
	return user
}

// 绑定输入参数
func BindInput(c *gin.Context, input interface{}) erron.E {
	if err := checkInput(input); err == nil {
		if err := c.ShouldBind(input); err != nil {
			return erron.Fail(erron.ErrRequestParams, err.Error())
		}
		// 如果实现了 params 接口就验证参数
		if obj, ok := input.(paramsCheck); ok {
			if err := obj.Check(); err != nil {
				return erron.Fail(erron.ErrRequestParams, err.Error())
			}
		}
	}
	return nil
}

// 展示的分页
func GetPager(c *gin.Context) *Pager {
	pager := &Pager{}
	c.ShouldBind(pager)
	pager.Secure()
	return pager
}

func setResponse(c *gin.Context, resp *response) {
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
