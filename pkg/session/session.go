package session

import (
	"Blog/pkg/config"
	"Blog/pkg/logger"
	"github.com/gorilla/sessions"
	"net/http"
)


// Store gorilla sessions 的存储库
var Store = sessions.NewCookieStore([]byte(config.GetString("app.key")))

//当前绘画  6212262507004524122    320981199310025717  13265735835
//霍建国  15571768956   中国建设银行  6217002830008056030   420500195901240617
//当前会话
var Session *sessions.Session //20210413543660566761842719182848

//用以获取会话
var Request *http.Request

//用于写入会话
var Response http.ResponseWriter

//初始化会话，在中间件中调用
func StartSession(w http.ResponseWriter,r *http.Request) {
	var err error

	//Store.Get() 的第二个参数是 Cookie的名称
	//gorilla/sessions 支持多会话，本项目只使用单一会话
	Session, err = Store.Get(r, config.GetString("session.session_name"))
	logger.LogError(err)

	Request = r
	Response = w
}

//写入键值对应的的会话数据
func Put(key string , value interface{})  {
	Session.Values[key] = value
	Save()
}

//获取会话数据，获取数据时请做类型检测
func Get(key string) interface{}  {
	return Session.Values[key]
}

//删除某个会话项
func Forget(key string) {
	delete(Session.Values,key)
	Save()
}

//删除当前会话
func Flush()  {
	Session.Options.MaxAge = -1
}

//保持会话
func Save()  {
	//非HTTPS的连接不能使用Secure 和 HttpOnly,浏览器会报错
	//Session.Options.Secure = true
	//Session.Options.HttpOnly = true

	err := Session.Save(Request, Response)
	logger.LogError(err)
}


