package routes

import (
	"Blog/app/http/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterWebRoutes(r *mux.Router)  {
	//静态页面
	pc := new(controllers.PagesController)
	r.HandleFunc("/",pc.Home).Methods("GET").Name("home")
	r.HandleFunc("/about",pc.About).Methods("GET").Name("about")
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)
}
