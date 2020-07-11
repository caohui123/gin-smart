package param

type SampleUser struct {
	Name string `json:"name" form:"name" binding:"name"`
}

type SampleUserResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}

type LoginRequest struct {
	Mobile string `form:"mobile" binding:"required"`
	Pwd    string `form:"pwd" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// form, binding 是gin的验证规则
// more validate rules refer: https://godoc.org/gopkg.in/go-playground/validator.v8
type UserListRequest struct {
	// Mobile string `form:"mobile" binding:"required"`
	Mobile string `form:"mobile" json:"mobile"`
	// Page   int64  `form:"page" binding:"required,gt=0"`
}

func (u *UserListRequest) Check() error {
	return nil
}

type SampleUserListResponse []UserItem

type UserDetailRequest struct {
	ID int64 `json:"id" form:"id" binding:"required"`
}

type UserDetailResponse struct {
	Id     uint   `json:"id"`
	Mobile string `json:"mobile"`
	Name   string `json:"name"`
}

//
type UserAddRequest struct {
	Mobile string `form:"mobile" binding:"required,len=11"`
	Pwd    string `form:"pwd" binding:"required,min=6,max=32"` // 6<= len(pwd) <=32
	Name   string `form:"name" binding:""`
}

type UserAddResponse struct {
	Id int64 `json:"id"`
}

type UserItem struct {
	Id     int64  `json:"id"`
	Mobile string `json:"mobile"`
	Name   string `json:"name"`
}

type UserModifyPwdRequest struct {
	Pwd    string `json:"pwd" form:"pwd" binding:"required"`
	PwdNew string `json:"pwd_new" form:"pwd_new" binding:"required"`
}
