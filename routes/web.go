package routes

import (
	"Blog/app/http/controllers"
	"Blog/app/http/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterWebRoutes(r *mux.Router)  {
	//静态页面
	pc := new(controllers.PagesController)
	r.HandleFunc("/about",pc.About).Methods("GET").Name("about")
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)

	//文章相关页面
	ac := new(controllers.ArticlesController)
	r.HandleFunc("/",ac.Index).Methods("GET").Name("home")
	r.HandleFunc("/articles/{id:[0-9]+}",ac.Show).Methods("GET").Name("articles.show")
	r.HandleFunc("/articles", ac.Index).Methods("GET").Name("articles.index")
	r.HandleFunc("/articles",middlewares.Auth(ac.Store)).Methods("POST").Name("articles.store")
	r.HandleFunc("/articles/create",middlewares.Auth(ac.Create)).Methods("GET").Name("articles.create")
	r.HandleFunc("/articles/{id:[0-9]+}/edit",middlewares.Auth(ac.Edit)).Methods("GET").Name("edit")
	r.HandleFunc("/articles/{id:[0-9]+}",middlewares.Auth(ac.Update)).Methods("POST").Name("articles.update")
	r.HandleFunc("/articles/{id:[0-9]+}/delete",middlewares.Auth(ac.Delete)).Methods("POST").Name("articles.delete")

	//用户认证
	auc := new(controllers.AuthController)
	r.HandleFunc("/auth/register",middlewares.Guest(auc.Register)).Methods("GET").Name("auth.register")
	r.HandleFunc("/auth/do-register",middlewares.Guest(auc.DoRegister)).Methods("POST").Name("auth.doregister")
	r.HandleFunc("/auth/login",middlewares.Guest(auc.Login)).Methods("GET").Name("auth.login")
	r.HandleFunc("/auth/dologin",middlewares.Guest(auc.DoLogin)).Methods("POST").Name("auth.dologin")
	r.HandleFunc("/auth/logout",middlewares.Guest(auc.Logout)).Methods("POST").Name("auth.logout")


	//静态资源
	r.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	r.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))

	//中间件，强制内容类型为HTML
	//r.Use(middlewares.ForceHTML)
	r.Use(middlewares.StartSession)
}
