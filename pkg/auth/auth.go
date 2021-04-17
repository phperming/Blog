package auth

import (
	"Blog/app/models/user"
	"Blog/pkg/session"
	"errors"
	"gorm.io/gorm"
)

func _GetUID() string {
	_uid := session.Get("uid")

	uid,ok := _uid.(string)

	if ok && len(uid) > 0 {
		return uid
	}

	return ""
}

//获取用户登录信息
func User() user.User {
	uid := _GetUID()

	if len(uid) > 0 {
		_user,err := user.Get(uid)
		if err == nil {
			return _user
		}
	}

	return user.User{}
}

//尝试登录
func Attempt(email string,password string) error {
	//1.根据email获取用户
	_user,err := user.GetByEmail(email)

	// 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("账号不存在或者密码错误")
		} else {
			return errors.New("内部错误，请稍后尝试")
		}
	}

	//匹配密码
	if !_user.ComparePassword(password) {
		return errors.New("账号不存在或者密码错误")
	}

	//登录用户，保存会话
	session.Put("uid",_user.GetStringID())

	return nil
}

//登录指定用户
func Login(_user user.User)  {
	session.Put("uid",_user.GetStringID())
}

//用户推出登录
func Logout()  {
	session.Forget("uid")
}

//检测是否登录
func Check() bool {
	return len(_GetUID()) > 0
}
