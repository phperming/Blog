package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprint(w,"<h1>Hello World</h>")
}

func main()  {
	http.HandleFunc("/",handlerFunc)
	http.ListenAndServe(":8088",nil)
}
