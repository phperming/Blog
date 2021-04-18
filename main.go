package main

import (
	"Blog/app/http/middlewares"
	"Blog/bootstrap"
	"Blog/config"
	c "Blog/pkg/config"
	"Blog/pkg/database"
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
)

var router *mux.Router
var db *sql.DB

func init()  {
	//初始化配置信息
	config.Initialize()
}

func main()  {
	database.Initialize()
	db = database.DB

	bootstrap.SetupDB()
	router =  bootstrap.SetupRoute()

	http.ListenAndServe(":" + c.GetString("app.port") ,
		middlewares.RemoveTrailingSlash(router))
}
