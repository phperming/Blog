package controllers

import (
	"Blog/app/models/user"
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

	//初始化变量
	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	//1.表单验证
	//2.验证通过  插入数据库 跳转到首页
	_user := user.User{
		Name: name,
		Email: email,
		Password: password,
	}
	_user.Create()

	if _user.ID > 0 {
		fmt.Fprint(w,"插入成功 ，ID为"+_user.GetStringID())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w,"创建用户失败。请联系管理员")
	}
	//3.验证不通过	重新显示表单
}