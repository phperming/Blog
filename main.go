package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

var router = mux.NewRouter()

func homeHandler(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprint(w,"<h1>Hello,欢迎来到GoBlog!</h1>")
}

func notFoundHandler(w http.ResponseWriter,r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w,"<h1>您要访问的页面不存在:(</h1> <p> 如有疑惑请联系我们</p>")

}

func aboutHandler(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprint(w,"此博客是用来记录变成笔记，如有反馈和建议请联系<a href=\"mailto:michel@163.com\">Michel@163.com</a>")
}

func articlesShowHandler(w http.ResponseWriter,r *http.Request)  {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w,"文章 ID:" + id)
}

func articlesIndexHandler(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprint(w,"访问文章列表")
}

func articlesStoreHandler(w http.ResponseWriter,r *http.Request) {
	fmt.Fprint(w,"创建新的文章")
}

func forceHTMLMiddleware(next http.Handler)http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//1.设置标头
		w.Header().Set("Content-Type","text/html;charset=utf-8")
		//2.继续处理请求
		next.ServeHTTP(w,r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//除首页以外，移除所有请求路径后面的/
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path,"/")
		}

		//将请求传递下去
		next.ServeHTTP(w,r)
	})
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<title>创建文章 —— 我的技术博客</title>
	</head>
	<body>
		<form action="%s" method="post">
			<p><input type="text" name="title"></p>
			<p><textarea name="body" cols="30" rows="10"></textarea></p>
			<p><button type="submit">提交</button></p>
		</form>
	</body>
	</html>`

	storeURL ,_:= router.Get("articles.store").URL()
	fmt.Fprintf(w,html,storeURL)
}


func main()  {
	router.HandleFunc("/",homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about",aboutHandler).Methods("GET").Name("about")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")

	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles",articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create",articlesCreateHandler).Methods("GET").Name("articles.create")


	//自定义404页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	//中间件 强制内容类型为HTML
	router.Use(forceHTMLMiddleware)

	//通过命名路由获取URL示例
	homeURL ,_:= router.Get("home").URL()
	fmt.Println(homeURL)
	articleURL ,_:= router.Get("articles.show").URL()
	fmt.Println(articleURL)

	http.ListenAndServe(":8088",removeTrailingSlash(router))
}
