package service

import (
	"errors"
	"fmt"

	"github.com/jangozw/gin-smart/pkg/util"
	"github.com/jangozw/go-api-facility/auth"

	"github.com/jangozw/gin-smart/pkg/app"
)

// 登陆业务实现
type JwtUserLogin struct {
	AccountID  string
	AccountPwd string
}

func (j *JwtUserLogin) Account() string {
	return j.AccountID
}

func (j *JwtUserLogin) Pwd() string {
	return j.AccountPwd
}

func (j *JwtUserLogin) JwtSecret() string {
	return app.Runner.Cfg.Encrypt.JwtSecret
}

func (j *JwtUserLogin) ExpireSeconds() int64 {
	return int64(app.Runner.Cfg.General.TokenExpireSeconds)
}

func (j *JwtUserLogin) Issuer() string {
	return "test"
}

func (j *JwtUserLogin) RefreshToken(accountID string, token auth.AccountToken) (err error) {
	// 存储 token
	userID := util.StrToInt64(accountID)
	return refreshUserToken(userID, string(token))
}

// 验证用户名，密码 根据自己的场景实现
func (j *JwtUserLogin) VerifyAccount(account string, pwd string) (*auth.AccountInfo, error) {
	userInfo, err := verifyAccount(account, pwd)
	if err != nil {
		return nil, err
	}
	return &auth.AccountInfo{
		AccountID: fmt.Sprintf("%d", userInfo.ID),
		AccountExtra: map[string]string{
			"name":   userInfo.Name,
			"mobile": userInfo.Mobile,
		},
	}, nil
}

// 验证token 业务实现
type JwtUserVerify struct {
	token string
}

func (j *JwtUserVerify) Token() string {
	return j.token
}

func (j *JwtUserVerify) SetToken(token string) {
	j.token = token
}

func (j *JwtUserVerify) JwtSecret() string {
	return app.Runner.Cfg.Encrypt.JwtSecret
}

// 可以不验证，直接返回nil, 或者从缓存，数据验证
func (j *JwtUserVerify) VerifyToken(accountID string, unbelievableAccountToken auth.AccountToken) (callback interface{}, err error) {
	if unbelievableAccountToken == "" {
		err = errors.New("token 为空")
		return
	}
	userID := util.StrToInt64(accountID)
	if userID <= 0 {
		err = errors.New("userID 为空")
		return
	}
	if userToken, e := redisGetLoginUser(userID); e != nil {
		err = e
		return
	} else if userToken == string(unbelievableAccountToken) {
		// 验证成功
		// callback = xxx
		callback = ""
		return callback, nil
	}
	return callback, errors.New("验证token失败")
}

/******************************************************************/
// 登陆举例
func ExampleLogin() {
	login := &JwtUserLogin{
		AccountID:  "13555555555",
		AccountPwd: "123456",
	}
	jwtTokenString, err := auth.Login(login)
	if err != nil {
		return
	}
	fmt.Println("login success, jwt-token:", jwtTokenString)
}

// 验证token 举例

func ExampleVerifyToken() {
	token := "a-s08238e209-3"
	verify := &JwtUserVerify{}
	jwtPayload, callback, err := auth.VerifyToken(verify, token)
	if err != nil {
		return
	}
	fmt.Println("verify success, info:", jwtPayload.AccountInfo, callback)
}
