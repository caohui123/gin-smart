package service

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"

	"github.com/jangozw/gin-smart/model"
	"github.com/jangozw/gin-smart/param"
	"github.com/jangozw/gin-smart/pkg/app"
)

const RedisKeyLoginUser = "login_user_token_"

// try login a user and return the token to client, client need to store the receive token and put it in http header before api request
func verifyAccount(account string, pwd string) (*model.SampleUser, error) {
	var user model.SampleUser
	if err := app.Runner.Db.Model(&user).Where("mobile=?", account).First(&user).Error; err != nil {
		return nil, err
	}
	if user.CheckPwd(pwd) != true {
		return nil, errors.New("invalid account or pwd")
	}
	return &user, nil
}

func AppLogout(userId int64) error {
	return app.Runner.Redis.DelKey(loginUserRedisKey(userId))
}

// 重新生成用户的token
func refreshUserToken(userID int64, token string) (err error) {
	// 存储 token
	var userToken model.SampleUserToken
	if err = app.Runner.Db.Model(&model.SampleUserToken{}).Where("user_id=?", userID).First(&userToken).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	expSeconds := int64(app.Runner.Cfg.General.TokenExpireSeconds)
	if userToken.ID == 0 {
		data := model.SampleUserToken{
			UserID:    userID,
			Token:     string(token),
			ExpiredAt: time.Now().Unix() + expSeconds,
		}
		err = app.Runner.Db.Create(&data).Error
	} else {
		err = app.Runner.Db.Model(&model.SampleUserToken{}).Where("id=?", userToken.ID).Update(map[string]interface{}{
			"expired_at": time.Now().Unix() + expSeconds,
			"token":      token,
		}).Error
	}
	if err != nil {
		return errors.New("create user token data failed")
	}
	if err = redisSetLoginUser(userID, string(token), expSeconds); err != nil {
		return err
	}
	return nil
}

// set user's token in expires
func redisSetLoginUser(userId int64, token string, exp int64) error {
	if err := app.Runner.Redis.SetKey(loginUserRedisKey(userId), token, exp); err != nil {
		return errors.New(fmt.Sprintf("redis set login user failed:%d, %s", userId, err.Error()))
	}
	return nil
}

//
func redisGetLoginUser(userId int64) (string, error) {
	return app.Runner.Redis.GetKey(loginUserRedisKey(userId))
}

func loginUserRedisKey(userId int64) string {
	return fmt.Sprintf("%s_%d", RedisKeyLoginUser, userId)
}

// 用户列表
func SampleGetUserList(search param.UserListRequest, pager app.Pager) (data param.SampleUserListResponse, err error) {
	var users []model.SampleUser
	query := app.Runner.Db.Model(&model.SampleUser{})
	if search.Mobile != "" {
		query = query.Where("mobile = ?", search.Mobile)
	}
	var total uint
	if err = query.Count(&total).Error; err != nil {
		return
	}
	if err = query.Limit(pager.Limit()).Offset(pager.Offset()).Find(&users).Error; err != nil {
		return
	}
	for _, u := range users {
		data = append(data, param.UserItem{
			Id:     u.ID,
			Mobile: u.Mobile,
			Name:   u.Name,
		})
	}
	// data.SetPager(search.Page, search.PageSize, total)
	return
}

func SampleGetUserByID(id uint) (*model.SampleUser, error) {
	var user model.SampleUser
	if err := app.Runner.Db.Where("id=?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
