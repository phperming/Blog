package main

import (
	"Blog/bootstrap"
	"Blog/pkg/database"
	"Blog/pkg/logger"
	"Blog/pkg/route"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

var router *mux.Router
var db *sql.DB

type ArticlesFormData struct {
	Title string
	Body string
	URL *url.URL
	Errors map[string]string
}

type Article struct {
	Title string
	Body string
	ID int64
}



func (a Article)Delete() (rowsAffected int64,err error)  {
	res, err := db.Exec("DELETE FROM articles WHERE id=" + strconv.FormatInt(a.ID, 10))
	if err != nil {
		return 0,err
	}

	//删除成功，跳转到文章详情页
	if n,_ := res.RowsAffected();n>0 {
		return n,nil
	}

	return 0,nil
}





func saveArticleToDB(title string,body string)(int64,error) {
	//变量初始化
	var (
		id int64
		err error
		rs sql.Result
		stmt *sql.Stmt
	)

	//1.获取一个prepare声明语句
	stmt, err = db.Prepare("INSERT INTO articles (title,body) VALUES (?,?)")
	//例行的错误检测
	if err != nil {
		return 0,err
	}

	//在此函数运行结束后关闭此语句，防止占用SQL连接
	defer stmt.Close()

	//执行请求，传参进入绑定内容
	rs, err = stmt.Exec(title, body)
	if err != nil {
		return  0,err
	}

	//成功的话会返回自增I
	if id,err = rs.LastInsertId(); id > 0 {
		return id,nil
	}

	return 0,err

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



func getArticleById(id string)(Article,error)  {
	article := Article{}
	query := "SELECT * FROM articles where id=?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article,err
}

func validateArticleFormData(title string,body string) map[string]string {
	errors := make(map[string]string)

	//验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于3-40之间"
	}

	//验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度不能小于10"
	}

	return errors
}


func articlesDeleteHandler(w http.ResponseWriter,r *http.Request) {
	//获取id
	id := route.GetRouterVariable("id", r)

	//读取对应的文章
	article, err := getArticleById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w,"404 文章未找到")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"500 服务器内部错误")
		}
	} else {
		//未出现错误，执行删除操作
		rowsAffected,err := article.Delete()

		//如果发生错误
		if err != nil {
			//应该SQL报错了
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"500 服务器内部错误")
		} else {
			//未发生错误
			if rowsAffected > 0 {
				//重定向到列表页
				indexURL ,_ := router.Get("articles.index").URL()
				http.Redirect(w,r,indexURL.String(),http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w,"404 文章未找到")
			}
		}
	}
}

func main()  {
	database.Initialize()
	db = database.DB

	bootstrap.SetupDB()
	router =  bootstrap.SetupRoute()



	router.HandleFunc("/articles/{id:[0-9]+}/delete",articlesDeleteHandler).Methods("POST").Name("articles.delete")

	//中间件 强制内容类型为HTML
	router.Use(forceHTMLMiddleware)

	//通过命名路由获取URL示例
	homeURL ,_:= router.Get("home").URL()
	fmt.Println(homeURL)
	articleURL ,_:= router.Get("articles.show").URL()
	fmt.Println(articleURL)

	http.ListenAndServe(":8088",removeTrailingSlash(router))
}
