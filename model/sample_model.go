package model

import (
	"errors"

	"github.com/jinzhu/gorm"

	"github.com/jangozw/gin-smart/pkg/app"
	"github.com/jangozw/gin-smart/pkg/util"
)

const (
	UserStatusNormal    = 1
	UserStatusForbidden = 2
)

// User 用户表
type SampleUser struct {
	gorm.Model
	Name     string // 姓名
	Mobile   string `gorm:"index"` // 手机号
	Password string // 密码
	Status   int8   // 状态
}

// UserToken 用户token表
type SampleUserToken struct {
	gorm.Model
	UserID    int64  `gorm:"index"` // 用户id
	Token     string // token
	ExpiredAt int64  // 过期时间
}

//
func SampleAddUser(name, mobile, pwd string) (SampleUser, error) {
	var total int
	if err := app.Db.Model(&SampleUser{}).Where("mobile=?", mobile).Count(&total).Error; err != nil {
		return SampleUser{}, err
	}
	if total > 0 {
		return SampleUser{}, errors.New("account already exist")
	}
	user := SampleUser{
		Name:     name,
		Mobile:   mobile,
		Password: SampleMakeUserPwd(pwd),
	}
	return user, app.Db.Create(&user).Error
}

func SampleFindUserByMobile(mobile string) (user SampleUser, err error) {
	if err = app.Db.Where("mobile=?", mobile).First(&user).Error; err != nil {
		return
	}
	return user, nil
}

func SampleMakeUserPwd(input string) string {
	return util.Sha256(input + "")
}

func (m *SampleUser) CheckPwd(input string) bool {
	return m.Password == util.Sha256(input+"")
}
