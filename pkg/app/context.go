package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	CtxKeyLoginUser = "ctx-login-user"
	CtxKeyNeedLogin = "ctx-need-login"
)

// const CtxKeyLoginUserID = "ctx-login-user-id"
const CtxKeyResponse = "ctx-response"

var (
	inputTypeErr  = errors.New("input 必须是一个结构体变量的地址")
	outputTypeErr = errors.New("output 必须是一个地址")
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

// 请求上下文, 基于 gin.Context
type baseContext struct {
	*gin.Context
}

// 用户api 请求
type Context struct {
	baseContext
	apiSetting ctxApiSetting
	Pager      Pager
	LoginUser  LoginUser
}

type ctxApiSetting struct {
	bindInput      interface{}
	bindOutput     interface{}
	route          Route
	outputPager    bool
	outputWithCode bool
	recordsCount   uint64
}

// 分页请求展示结果
type pagerResponse struct {
	Pager pagerResp   `json:"pager"`
	List  interface{} `json:"list"`
}

// 用户api类型
// type APIHandlerFunc func(c *Context)

func APIHandler(route Route) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 新建ctx
		ctx, err := newCtx(c, route)
		if err != nil {
			c.JSON(http.StatusOK, ResponseFail(err))
			return
		}
		// 支持两种handler形式
		// 1, 通过定义func 直接处理
		if route.HandlerFunc != nil {
			if err := route.HandlerFunc(ctx); err != nil {
				c.JSON(http.StatusOK, ResponseFail(err))
				return
			}
			// 2, 通过指定的对象实例处理请求
		} else if route.Handler != nil {
			newInstance := route.Handler()
			if newInstance == nil {
				c.JSON(http.StatusOK, ResponseFailByCode(errCodeNotSetRouteHandler))
				return
			}
			if err := newInstance.Prepare(ctx); err != nil {
				c.JSON(http.StatusOK, ResponseFail(err))
				return
			}
			if err := newInstance.Handler(); err != nil {
				c.JSON(http.StatusOK, ResponseFail(err))
				return
			}
			// 没有找到handler
		} else {
			c.JSON(http.StatusOK, ResponseFailByCode(errCodeNotSetRouteHandler))
			return
		}
		// 到这里都是请求成功的
		// 根据情形调整output
		ctx.adjustOutput()
		var result interface{}
		if ctx.apiSetting.outputWithCode != true {
			result = ctx.apiSetting.bindOutput
		} else {
			result = ctx.Success(ctx.apiSetting.bindOutput)
		}
		// 设置输出结果用于日志记录
		c.Set(CtxKeyResponse, result)
		c.JSON(http.StatusOK, result)
		return
	}
}

// 新 context
func newCtx(c *gin.Context, ro Route) (*Context, Err) {
	/******* 对于需要登陆情况检查登陆 ******/
	var loginUser LoginUser
	if c.GetBool(CtxKeyNeedLogin) == true {
		if obj, exists := c.Get(CtxKeyLoginUser); exists {
			loginUser = obj.(LoginUser)
		}
		if loginUser.ID == 0 {
			return nil, Error(errCodeInvalidLoginUser, "用户id无效0")
		}
	}
	// 通用参数
	page := StringNumber(c.Query(URLParamPage))
	pageSize := StringNumber(c.Query(URLParamPageSize))
	pager := NewPageRequest(page.Uint(), pageSize.Uint())

	// 根据路由自动判断是否是分页请求，如果是则返回结果自动加上pager
	var outputPager bool
	if ro.Method == "GET" && strings.HasSuffix(ro.Path, "list") {
		outputPager = true
	}
	ctx := &Context{
		baseContext: baseContext{c},
		apiSetting: ctxApiSetting{
			route:          ro,
			bindInput:      nil,
			bindOutput:     nil,
			outputPager:    outputPager,
			recordsCount:   0,
			outputWithCode: true,
		},
		Pager:     pager,
		LoginUser: loginUser,
	}
	return ctx, nil
}

func (c *Context) adjustOutput() {
	// 如果是分页请求,改变下返回结构
	if c.apiSetting.outputPager == true {
		c.apiSetting.bindOutput = pagerResponse{
			Pager: pagerResp{
				Page:     c.Pager.Page,
				PageSize: c.Pager.PageSize,
				Total:    c.apiSetting.recordsCount,
			},
			List: c.apiSetting.bindOutput,
		}
	}
}

// 展示分页
func (c *Context) PagerOn(count uint64) {
	c.apiSetting.outputPager = true
	c.apiSetting.recordsCount = count
}

func (c *Context) PagerOff() {
	c.apiSetting.outputPager = false
}

// 绑定输入参数
func (c *Context) BindInput(input interface{}) Err {
	if err := checkInput(input); err == nil {
		if err := c.ShouldBind(input); err != nil {
			return Error(errCodeInvalidRequestParam, err)
		}
		// 如果实现了 params 接口就验证参数
		if obj, ok := input.(paramsCheck); ok {
			if err := obj.Check(); err != nil {
				return Error(errCodeInvalidRequestParam, err)
			}
		}
		c.apiSetting.bindInput = input
	}
	return nil
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

// 请求入参出参绑定
func (c *Context) MustBinds(input interface{}, output interface{}) error {
	if err := c.BindInput(input); err != nil {
		return err
	}
	if err := c.bindOutput(output); err != nil {
		return err
	}
	return nil
}

// 预先绑定输出结构
func (c *Context) bindOutput(output interface{}) error {
	// 非生产环境判断
	if gin.Mode() == gin.DebugMode {
		rv := reflect.ValueOf(output)
		if rv.Kind() != reflect.Ptr || rv.IsNil() || !rv.IsValid() {
			return outputTypeErr
		}
	}
	c.apiSetting.bindOutput = output
	return nil
}

// 设置输出结果
func (c *Context) Output(v interface{}) {
	c.apiSetting.bindOutput = v
}

// 不带外包装结构code，msg 直接输出
func (c *Context) OutputWithoutWrapping(v interface{}) {
	c.apiSetting.outputWithCode = false
	c.apiSetting.bindOutput = v
}

// 返回成功， 有数据
func (c *Context) Success(data interface{}) *Response {
	return c.response(nil, data)
}

// 返回成功，不带数据，通用型
func (c *Context) SuccessSimple() *Response {
	return c.response(nil, struct{}{})
}

// 返回错误， 通用型
func (c *Context) Fail(msg string, args ...interface{}) *Response {
	err := Error(errCodeFail, fmt.Sprintf(msg, args...))
	return c.response(err, struct{}{})
}

// 返回错误， 带错误信息
func (c *Context) FailWithCode(err Err) *Response {
	return c.response(err, struct{}{})
}

func (c *Context) response(err Err, data interface{}) *Response {
	return ResponseWithCode(err, data)
}
