package main

import (
	"fmt"
	"net/http"
	"strings"
)

func defaultHandler(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","text/html;charset=utf-8")
	if r.URL.Path == "/" {
		fmt.Fprint(w,"<h1>Hello ,欢迎来到 GOBlog</h1>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w,"<h1>您要访问的页面不存在:(</h1> <p> 如有疑惑请联系我们</p>")
	}
}

func aboutHandler(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","text/html;charset=utf-8")

	fmt.Fprint(w,"此博客是用来记录变成笔记，如有反馈和建议请联系<a href=\"mailto:michel@163.com\">Michel@163.com</a>")

}

func main()  {
	router := http.NewServeMux()

	router.HandleFunc("/",defaultHandler)
	router.HandleFunc("/about",aboutHandler)
	router.HandleFunc("/articles/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.SplitN(r.URL.Path,"/",3)[2]
		fmt.Fprint(w,"文章 ID:" + id)
	})

	router.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fmt.Fprint(w,"文章列表")
		case "POST":
			fmt.Fprint(w,"创建新的文章")

		}
	})
	http.ListenAndServe(":8088",router)
}
