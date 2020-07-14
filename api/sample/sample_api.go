package sample

import (
	"github.com/jangozw/gin-smart/erro"
	"github.com/jangozw/gin-smart/model"
	"github.com/jangozw/gin-smart/param"
	"github.com/jangozw/gin-smart/pkg/app"
	"github.com/jangozw/gin-smart/pkg/auth"
	"github.com/jangozw/gin-smart/service"
)

// login api
func Login(c *app.Context) (interface{}, erro.E) {
	var input param.LoginRequest
	if err := c.ShouldBind(&input); err != nil {
		return nil, erro.Fail(erro.ErrRequestParams, err.Error())
	}
	user, err := model.FindUserByMobile(input.Mobile)
	if err != nil {
		return nil, erro.FailBy(err)
	}
	if !user.CheckPwd(input.Pwd) {
		return nil, erro.New("invalid account or pwd")
	}
	// token 携带的user 信息根据业务情况设置
	tokenPayload := app.TokenPayload{UserID: user.ID}
	token, err := auth.GenerateJwtToken(app.Cfg.General.JwtSecret, app.Cfg.General.TokenExpire, tokenPayload)
	if err != nil {
		return nil, erro.FailBy(err)
	}
	output := param.LoginResponse{Token: token}
	return output, nil
}

// 带有分页的列表
func UserList(c *app.Context) (interface{}, erro.E) {
	input := param.UserListRequest{}
	if err := c.BindInput(&input); err != nil {
		return nil, err
	}
	pager := c.GetPager()
	list, err := service.SampleGetUserList(input, pager)
	if err != nil {
		return nil, erro.FailBy(err)
	}
	return app.PagerResponse(pager, list), nil
}

// logout api
func Logout(c *app.Context) (interface{}, erro.E) {
	userId := c.MustLoginUser().ID
	if err := service.AppLogout(int64(userId)); err != nil {
		return nil, erro.FailBy(err)
	}
	return nil, nil
}

// 添加用户
func AddUser(c *app.Context) (interface{}, erro.E) {
	var input param.UserAddRequest
	if err := c.BindInput(&input); err != nil {
		return nil, erro.FailBy(err)
	}
	if user, err := model.AddUser(input.Name, input.Mobile, input.Pwd); err != nil {
		return nil, erro.FailBy(err)
	} else {
		data := param.SampleUserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Mobile: user.Mobile,
		}
		return data, nil
	}
}

// 用户详情
func UserDetail(c *app.Context) (interface{}, erro.E) {
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
	return output, nil
}

// 修改的自己的密码
func UserChangePwd(c *app.Context) (interface{}, erro.E) {
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
func UserListTest(c *app.Context) (interface{}, erro.E) {
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
