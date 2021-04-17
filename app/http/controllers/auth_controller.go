package controllers

import (
	"Blog/app/models/user"
	"Blog/app/requests"
	"Blog/pkg/auth"
	"Blog/pkg/flash"
	"Blog/pkg/session"
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
			//注册成功， 登录用户并跳转到首页
			flash.Success("恭喜您注册成功")
			auth.Login(_user)
			http.Redirect(w,r,"/",http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"创建用户失败。请联系管理员")
		}
	}

}

func (*AuthController)Login(w http.ResponseWriter,r *http.Request)  {
	session.Put("uid","1")
	view.RendSimple(w,view.D{},"auth.login")
}

func (*AuthController)DoLogin(w http.ResponseWriter,r *http.Request)  {
	//初始化表单数据
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	//尝试登录
	if err := auth.Attempt(email,password); err == nil {
		//登录成功
		flash.Success("欢迎回来！")
		http.Redirect(w,r,"/",http.StatusFound)
	} else {
		//失败 显示错误信息
		view.RendSimple(w,view.D{
			"Error" : err.Error(),
			"Email" : email,
			"Password" : password,
		},"auth.login")

	}

}

//退出登录z
func (*AuthController)Logout(w http.ResponseWriter,r * http.Request)  {
	auth.Logout()
	flash.Success("您已退出登录")
	http.Redirect(w,r,"/",http.StatusFound)
}