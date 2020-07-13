package sample

import (
	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/erro"
	"github.com/jangozw/gin-smart/model"
	"github.com/jangozw/gin-smart/param"
	"github.com/jangozw/gin-smart/pkg/app"
	"github.com/jangozw/gin-smart/pkg/auth"
	"github.com/jangozw/gin-smart/service"
)

// login api

func SampleLogin(c *gin.Context) (interface{}, erro.E) {
	var p param.LoginRequest
	if err := c.ShouldBind(&p); err != nil {
		return nil, erro.Fail(erro.ErrRequestParams)
	}
	// token 携带的user 信息根据业务情况设置
	user := app.LoginUser{ID: 1}
	jwtToken, err := auth.GenerateJwtToken(app.Cfg.General.JwtSecret, app.Cfg.General.TokenExpire, user)
	if err != nil {
		return nil, erro.FailBy(err)
	}
	output := param.LoginResponse{Token: jwtToken}
	return output, nil
}

func Sample(c *app.Context) (interface{}, erro.E) {
	return "sample ok", nil
}

func List(c *app.Context) (interface{}, erro.E) {

	// 请求的分页信息
	/*
		page := ctx.Pager().Page
		pageSize := ctx.Pager().PageSize
		offset := ctx.Pager().Offset()
		limit := ctx.Pager().Limit()
	*/

	type user struct {
		Name string `json:"name"`
		Age  int8   `json:"age"`
	}
	list := make([]user, 0)
	list = append(list, user{
		Name: "Dog",
		Age:  2,
	})
	list = append(list, user{
		Name: "Cat",
		Age:  3,
	})

	var total uint = 2
	// 设置数据总数， 自动调整输出结构为分页结构
	c.SetPager(total)
	return list, nil
}

// logout api
func SampleLogout(c *app.Context) (interface{}, erro.E) {
	userId := c.MustLoginUser().ID
	if err := service.AppLogout(int64(userId)); err != nil {
		return nil, erro.FailBy(err)
	}
	return nil, nil
}

// 添加用户
func SampleAddUser(c *app.Context) (interface{}, erro.E) {
	var input param.UserAddRequest
	if err := c.BindInput(&input); err != nil {
		return nil, erro.FailCn()
	}
	if user, err := model.SampleAddUser(input.Name, input.Mobile, input.Pwd); err != nil {
		return nil, erro.FailCn()
	} else {
		data := param.SampleUserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Mobile: user.Mobile,
		}
		return data, nil
	}
}

// 用户列表，有分页
func SampleUserList(c *app.Context) (interface{}, erro.E) {
	// 校验请求参数, 校验规则定义在params.SearchUserList{}的tag里
	search := param.UserListRequest{}
	if err := c.BindInput(&search); err != nil {
		return nil, erro.FailCn()
	}
	// 校验参数成功后自动赋值给结构体
	users, err := service.SampleGetUserList(search, *c.Pager())
	if err != nil {
		return nil, erro.FailBy(err)
	}
	userList := make(param.SampleUserListResponse, len(users))
	for i, v := range users {
		userList[i].Id = v.Id
		userList[i].Name = v.Name
		userList[i].Mobile = v.Mobile
	}
	return userList, nil
}

// 用户详情
func SampleUserDetail(c *app.Context) (interface{}, erro.E) {
	var input param.UserDetailRequest
	if err := c.BindInput(&input); err != nil {
		return nil, erro.FailBy(err)
	}
	user, err := service.SampleGetUserByID(input.ID)
	if err != nil {
		return nil, erro.FailBy(err)
	}
	output := &param.UserDetailResponse{
		Id:     user.ID,
		Mobile: user.Mobile,
		Name:   user.Name,
	}
	// c.OutputWithoutWrapping(output)
	return output, nil
}

// 修改的自己的密码
func SampleUserModifyPwd(c *app.Context) (interface{}, erro.E) {
	var input param.UserModifyPwdRequest
	if err := c.BindInput(&input); err != nil {
		return nil, erro.FailCn()
	}
	// 当前登陆用户
	loginUid := c.MustLoginUser().ID
	_, err := service.SampleGetUserByID(loginUid)
	if err != nil {
		return nil, erro.FailBy(err)
	}
	// 修改密码
	// 。。。
	return nil, nil
}

// 用户列表
func SampleUser(c *app.Context) (interface{}, erro.E) {
	// 验证请求参数
	input := param.SampleUser{}
	if err := c.ShouldBind(&input); err != nil {
		return nil, erro.FailBy(err)
	}
	// 返回结果
	output := param.SampleUserResponse{
		Name:   "sample",
		Mobile: "135000000000",
	}
	return output, nil
}
