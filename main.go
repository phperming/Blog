package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func homeHandler(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","text/html;charset=utf-8")
	fmt.Fprint(w,"<h1>Hello,欢迎来到GoBlog!</h1>")
}

func notFoundHandler(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","text/html;charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w,"<h1>您要访问的页面不存在:(</h1> <p> 如有疑惑请联系我们</p>")

}

func aboutHandler(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","text/html;charset=utf-8")

	fmt.Fprint(w,"此博客是用来记录变成笔记，如有反馈和建议请联系<a href=\"mailto:michel@163.com\">Michel@163.com</a>")
}

func articlesShowHandler(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","text/html;charset=utf-8")
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

func main()  {
	router := mux.NewRouter()

	router.HandleFunc("/",homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about",aboutHandler).Methods("GET").Name("about")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")

	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles",articlesStoreHandler).Methods("POST").Name("articles.store")

	//自定义404页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	//通过命名路由获取URL示例
	homeURL ,_:= router.Get("home").URL()
	fmt.Println(homeURL)
	articleURL ,_:= router.Get("articles.show").URL()
	fmt.Println(articleURL)

	http.ListenAndServe(":8088",router)
}
