package controllers

import (
	"Blog/app/models/user"
	"Blog/app/requests"
	"Blog/pkg/view"
	"fmt"
	"net/http"
)

type AuthController struct {

}


func (*AuthController)Register(w http.ResponseWriter,r *http.Request)  {
	view.Render(w,view.D{},"auth.register")
}

func (*AuthController)DoRegister(w http.ResponseWriter,r *http.Request)  {

	//初始化数据
	_user := user.User{
		Name: r.PostFormValue("name"),
		Email: r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
		PasswordConfirm : r.PostFormValue("password_confirm"),
	}

	errs := requests.ValidateRegistrationForm(_user)

	if len(errs) > 0 {
		//有错误发生，打印数据
		view.RendSimple(w,view.D{
			"Errors" : errs,
			"User" : _user,
		},"auth.register")
	} else {
		//2.验证通过  插入数据库 跳转到首页

		_user.Create()

		if _user.ID > 0 {
			fmt.Fprint(w,"插入成功 ，ID为"+_user.GetStringID())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"创建用户失败。请联系管理员")
		}
	}

}