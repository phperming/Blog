package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var router = mux.NewRouter()
var db  *sql.DB


func initDB()  {
	var err error
	config := mysql.Config{

		User : "root",
		Passwd : "root",
		Addr: "127.0.0.1:3306",
		Net: "tcp",
		DBName: "goblog",
		AllowNativePasswords: true,
	}

	//准备数据库连接池
	db,err = sql.Open("mysql",config.FormatDSN())

	checkError(err)
	//设置最大连接数
	db.SetMaxOpenConns(25)
	//设置最大空闲数
	db.SetMaxIdleConns(25)
	//设置每个连接的过期时间
	db.SetConnMaxLifetime(5 * time.Minute)

	//尝试连接，失败会报错
	err = db.Ping()

	checkError(err)
}

func checkError(err error)  {
	if err != nil {
		log.Fatal(err)
	}
}

func createTables()  {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
		id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
		title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
		body longtext COLLATE utf8mb4_unicode_ci
	);
	`

	_, err := db.Exec(createArticlesSQL)
	checkError(err)

}

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
	//获取URL参数
	vars := mux.Vars(r)
	id := vars["id"]

	//读取对应的文章数据
	article := Article{}
	query := "SELECT * FROM articles WHERE id=?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	if err != nil {
		if err == sql.ErrNoRows {
			//数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w,"404,文章未找到")
		} else {
			//数据库错误
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"500,内部服务器错误")
		}
	} else {
		//文章读取成功，显示文章
		tmpl, err := template.ParseFiles("resource/views/articles/show.gohtml")
		checkError(err)
		tmpl.Execute(w,article)
	}

}

func articlesIndexHandler(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprint(w,"访问文章列表")
}

func articlesStoreHandler(w http.ResponseWriter,r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := make(map[string]string)

	//验证title
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于3-30"
	}

	//验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于10"
	}
	
	//检查是否有错误
	if len(errors) == 0 {
		lastInsertId,err := saveArticleToDB(title,body)
		if lastInsertId > 0 {
			fmt.Println("出插入成功，ID为"+strconv.FormatInt(lastInsertId,10))
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"500 服务器内部错误")
		}
	} else {
		storeURL ,_ := router.Get("articles.store").URL()
		data := ArticlesFormData{
			Title : title,
			Body: body,
			URL : storeURL,
			Errors: errors,
		}

		tmpl, err := template.ParseFiles("resource/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}
		tmpl.Execute(w,data)
	}
	
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

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	storeURL ,_:= router.Get("articles.store").URL()
	data := ArticlesFormData{
		Title: "",
		Body: "",
		URL : storeURL,
		Errors: nil,
	}
	tmpl, err := template.ParseFiles("resource/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w,data)
}


func main()  {
	initDB()
	createTables()
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
