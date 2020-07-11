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

func SampleLogin(c *gin.Context) {
	ctx := app.Ctx(c)
	var p param.LoginRequest
	if err := ctx.ShouldBind(&p); err != nil {
		ctx.Fail(erro.Info(erro.ErrRequestParams))
		return
	}
	// token 携带的user 信息根据业务情况设置
	user := app.LoginUser{ID: 1}
	jwtToken, err := auth.GenerateJwtToken(app.Cfg.General.JwtSecret, app.Cfg.General.TokenExpire, user)
	if err != nil {
		ctx.Fail(erro.Info(erro.Failed))
		return
	}
	output := param.LoginResponse{Token: string(jwtToken)}
	ctx.Success(output)
}

func Sample(c *gin.Context) {
	ctx := app.Ctx(c)
	ctx.Success("ok")
}

func List(c *gin.Context) {
	ctx := app.Ctx(c)

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
	ctx.SetPager(total)
	ctx.Success(list)
}

// logout api
func SampleLogout(c *gin.Context) {
	ctx := app.Ctx(c)
	userId := ctx.MustLoginUser().ID
	if err := service.AppLogout(int64(userId)); err != nil {
		ctx.Fail(erro.Info(erro.Failed))
		return
	}
	ctx.Success(nil)
}

// 添加用户
func SampleAddUser(c *gin.Context) {
	ctx := app.Ctx(c)
	var input param.UserAddRequest
	if err := ctx.BindInput(&input); err != nil {
		ctx.Fail(erro.Info(erro.Failed))
		return
	}
	if user, err := model.SampleAddUser(input.Name, input.Mobile, input.Pwd); err != nil {
		ctx.Fail(erro.Info(erro.Failed))
		return
	} else {
		ctx.Success(param.SampleUserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Mobile: user.Mobile,
		})
	}
}

// 用户列表，有分页
func SampleUserList(c *gin.Context) {
	ctx := app.Ctx(c)
	// 校验请求参数, 校验规则定义在params.SearchUserList{}的tag里
	search := param.UserListRequest{}
	if err := ctx.BindInput(&search); err != nil {
		ctx.Fail(erro.Info(erro.Failed))
		return
	}
	// 校验参数成功后自动赋值给结构体
	if users, err := service.SampleGetUserList(search, *ctx.Pager()); err != nil {
		ctx.Fail(erro.Info(erro.Failed))
		return
	} else {
		userList := make(param.SampleUserListResponse, len(users))
		for i, v := range users {
			userList[i].Id = v.Id
			userList[i].Name = v.Name
			userList[i].Mobile = v.Mobile
		}
		ctx.Success(userList)
	}
	return
}

// 用户详情
func SampleUserDetail(c *gin.Context) {
	ctx := app.Ctx(c)
	var input param.UserDetailRequest
	if err := ctx.BindInput(&input); err != nil {
		ctx.Fail(erro.Info(erro.Failed))
		return
	}
	user, err := service.SampleGetUserByID(input.ID)
	if err != nil {
		ctx.Fail(erro.Info(erro.Failed))
		return
	}
	output := &param.UserDetailResponse{
		Id:     user.ID,
		Mobile: user.Mobile,
		Name:   user.Name,
	}
	// c.OutputWithoutWrapping(output)
	ctx.Success(output)
	return
}

// 修改的自己的密码
func SampleUserModifyPwd(c *gin.Context) {
	ctx := app.Ctx(c)
	var input param.UserModifyPwdRequest
	if err := ctx.BindInput(&input); err != nil {
		ctx.Fail(erro.Info(erro.Failed))
		return
	}
	// 当前登陆用户
	loginUid := ctx.MustLoginUser().ID
	_, err := service.SampleGetUserByID(loginUid)
	if err != nil {
		ctx.Fail(erro.Info(erro.Failed))
		return
	}
	// 修改密码
	// 。。。
	ctx.Success(nil)
	return
}

// 用户列表
func SampleUser(c *gin.Context) {
	ctx := app.Ctx(c)

	// 验证请求参数
	input := param.SampleUser{}
	if err := ctx.BindInput(&input); err != nil {
		ctx.Fail(erro.Info(erro.Failed))
		return
	}
	// 返回结果
	output := param.SampleUserResponse{
		Name:   "sample",
		Mobile: "135000000000",
	}
	// 设置输出结果
	ctx.Success(output)
	return
}
