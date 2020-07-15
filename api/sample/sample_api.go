package sample

import (
	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/erron"
	"github.com/jangozw/gin-smart/model"
	"github.com/jangozw/gin-smart/param"
	"github.com/jangozw/gin-smart/pkg/app"
	"github.com/jangozw/gin-smart/pkg/auth"
	"github.com/jangozw/gin-smart/service"
)

// login api
func Login(c *gin.Context) (interface{}, erron.E) {
	var input param.LoginRequest
	if err := c.ShouldBind(&input); err != nil {
		return nil, erron.Fail(erron.ErrRequestParams, err.Error())
	}
	user, err := model.FindUserByMobile(input.Mobile)
	if err != nil {
		return nil, erron.FailBy(err)
	}
	if !user.CheckPwd(input.Pwd) {
		return nil, erron.New("invalid account or pwd")
	}
	// token 携带的user 信息根据业务情况设置
	tokenPayload := app.TokenPayload{UserID: user.ID}
	token, err := auth.GenerateJwtToken(app.Cfg.General.JwtSecret, app.Cfg.General.TokenExpire, tokenPayload)
	if err != nil {
		return nil, erron.FailBy(err)
	}
	output := param.LoginResponse{Token: token}
	return output, nil
}

// 带有分页的列表
func UserList(c *gin.Context) (interface{}, erron.E) {
	input := param.UserListRequest{}
	if err := c.ShouldBind(&input); err != nil {
		return nil, erron.FailBy(err)
	}
	pager := app.GetPager(c)
	list, err := service.SampleGetUserList(input, pager)
	if err != nil {
		return nil, erron.FailBy(err)
	}
	return app.PagerResponse(pager, list), nil
}

// logout api
func Logout(c *gin.Context) (interface{}, erron.E) {
	userId := app.MustGetLoginUser(c).ID
	if err := service.AppLogout(userId); err != nil {
		return nil, erron.FailBy(err)
	}
	return nil, nil
}

// 添加用户
func AddUser(c *gin.Context) (interface{}, erron.E) {
	var input param.UserAddRequest
	if err := c.ShouldBind(&input); err != nil {
		return nil, erron.FailBy(err)
	}
	if user, err := model.AddUser(input.Name, input.Mobile, input.Pwd); err != nil {
		return nil, erron.FailBy(err)
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
func UserDetail(c *gin.Context) (interface{}, erron.E) {
	var input param.UserDetailRequest
	if err := c.ShouldBind(&input); err != nil {
		return nil, erron.FailBy(err)
	}
	user, err := service.SampleGetUserByID(input.ID)
	if err != nil {
		return nil, erron.FailBy(err)
	}
	output := &param.UserDetailResponse{
		Id:     user.ID,
		Mobile: user.Mobile,
		Name:   user.Name,
	}
	return output, nil
}

// 修改的自己的密码
func UserChangePwd(c *gin.Context) (interface{}, erron.E) {
	var input param.UserModifyPwdRequest
	if err := app.BindInput(c, &input); err != nil {
		return nil, erron.FailCn()
	}
	// 当前登陆用户
	loginUid := app.MustGetLoginUser(c).ID
	_, err := service.SampleGetUserByID(loginUid)
	if err != nil {
		return nil, erron.FailBy(err)
	}
	// 修改密码
	// 。。。
	return nil, nil
}

// 用户列表
func UserListTest(c *gin.Context) (interface{}, erron.E) {
	// 验证请求参数
	input := param.SampleUser{}
	if err := c.ShouldBind(&input); err != nil {
		return nil, erron.FailBy(err)
	}
	// 返回结果
	output := param.SampleUserResponse{
		Name:   "sample",
		Mobile: "135000000000",
	}
	return output, nil
}
