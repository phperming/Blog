package bootstrap

import (
	"Blog/routes"
	"github.com/gorilla/mux"
)

//初始化路由
func SetupRoute() *mux.Router {
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)

	return router
}
