package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","text/html;charset=utf-8")
	if r.URL.Path == "/" {
		fmt.Fprint(w,"<h1>欢迎来到goBlog</h>")
	} else if r.URL.Path == "/about" {
		fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
			"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w,"你要找的页面不存在")
	}

}

func main()  {
	http.HandleFunc("/",handlerFunc)
	http.ListenAndServe(":8088",nil)
}