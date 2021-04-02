package main

import (
	"Blog/app/http/middlewares"
	"Blog/bootstrap"
	"Blog/pkg/database"
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
)

var router *mux.Router
var db *sql.DB

func main()  {
	database.Initialize()
	db = database.DB

	bootstrap.SetupDB()
	router =  bootstrap.SetupRoute()

	http.ListenAndServe(":8088",middlewares.RemoveTrailingSlash(router))
}
