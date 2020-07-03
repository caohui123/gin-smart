package sample

import (
	"github.com/jangozw/go-api-facility/auth"
	"github.com/jangozw/gin-smart/errcode"
	"github.com/jangozw/gin-smart/model"
	"github.com/jangozw/gin-smart/param"
	"github.com/jangozw/gin-smart/pkg/app"
	"github.com/jangozw/gin-smart/service"
)

//login api

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

//logout api
func SampleLogout(c *app.Context) app.Err {
	userId := c.LoginUser.ID
	if err := service.AppLogout(int64(userId)); err != nil {
		return app.Error(errcode.Failed, err.Error())
	}
	return nil
}

// 添加用户
func SampleAddUser(c *app.Context) app.Err {
	var input param.UserAddRequest
	if err := c.BindInput(&input); err != nil {
		return app.Error(errcode.ErrRequestParams, err)
	}
	if user, err := model.SampleAddUser(input.Name, input.Mobile, input.Pwd); err != nil {
		return app.Error(errcode.Failed, err)
	} else {
		c.Output(param.SampleUserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Mobile: user.Mobile,
		})
	}
	return nil
}

// 用户列表，有分页
func SampleUserList(c *app.Context) app.Err {
	//校验请求参数, 校验规则定义在params.SearchUserList{}的tag里
	search := param.UserListRequest{}
	if err := c.BindInput(&search); err != nil {
		return err
	}
	//校验参数成功后自动赋值给结构体
	if users, err := service.SampleGetUserList(search, c.Pager); err != nil {
		return app.Error(errcode.Failed, "get user list err: %s", err.Error())
	} else {
		userList := make(param.SampleUserListResponse, len(users))
		for i, v := range users {
			userList[i].Id = v.Id
			userList[i].Name = v.Name
			userList[i].Mobile = v.Mobile
		}
		c.Output(userList)
	}
	return nil
}

// 用户详情
func SampleUserDetail(c *app.Context) app.Err {
	var input param.UserDetailRequest
	if err := c.BindInput(&input); err != nil {
		return err
	}
	user, err := service.SampleGetUserByID(input.ID)
	if err != nil {
		return app.Error(errcode.Failed, "get user by id err: %s", err.Error())
	}
	output := &param.UserDetailResponse{
		Id:     user.ID,
		Mobile: user.Mobile,
		Name:   user.Name,
	}
	// c.OutputWithoutWrapping(output)
	c.Output(output)
	return nil
}

// 修改的自己的密码
func SampleUserModifyPwd(c *app.Context) app.Err {
	var input param.UserModifyPwdRequest
	if err := c.BindInput(&input); err != nil {
		return err
	}
	// 当前登陆用户
	loginUid := c.LoginUser.ID
	_, err := service.SampleGetUserByID(loginUid)
	if err != nil {
		return app.Error(errcode.Failed, "couldn't get user err: %s", err.Error())
	}
	// 修改密码
	// 。。。
	return nil

}

// 用户列表
func SampleUser(c *app.Context) app.Err {
	// 验证请求参数
	input := param.SampleUser{}
	if err := c.BindInput(&input); err != nil {
		return app.CodeErr(errcode.ErrRequestParams)
	}
	// 返回结果
	output := param.SampleUserResponse{
		Name:   "sample",
		Mobile: "135000000000",
	}
	// 设置输出结果
	c.Output(output)
	return nil
}
