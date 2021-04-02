package controllers

import (
	"Blog/app/models/article"
	"Blog/pkg/logger"
	"Blog/pkg/route"
	"Blog/pkg/types"
	"fmt"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"strconv"
	"unicode/utf8"
)

type ArticlesController struct {

}

type ArticlesFormData struct {
	Title string
	Body string
	URL string
	Errors map[string]string
}


func (*ArticlesController)Show(w http.ResponseWriter,r *http.Request)  {
	//获取URL参数
	id := route.GetRouterVariable("id",r)

	//读取对应的文章数据
	article,err := article.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w,"404,文章未找到")
		} else {
			//数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"500,内部服务器错误")
		}
	} else {
		//文章读取成功，显示文章
		tmpl,err := template.New("show.gohtml").Funcs(template.FuncMap{
			"RouteName2URL" : route.Name2URL,
			"Int64ToString" : types.Int64ToString,
		}).ParseFiles("resource/views/articles/show.gohtml")
		//tmpl, err := template.ParseFiles("resource/views/articles/show.gohtml")
		logger.LogError(err)
		tmpl.Execute(w,article)
	}
}

func (*ArticlesController)Index(w http.ResponseWriter,r *http.Request)  {
	articles,err := article.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w,"500 服务器内部错误")
	} else {
		tmpl,err := template.ParseFiles("resource/views/articles/index.gohtml")
		logger.LogError(err)
		tmpl.Execute(w,articles)
	}
}

func (*ArticlesController)Create(w http.ResponseWriter, r *http.Request) {
	storeURL := route.Name2URL("articles.store")
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

func (*ArticlesController)Store(w http.ResponseWriter,r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title,body)
	//检查是否有错误
	if len(errors) == 0 {
		_article := article.Article{
			Title: title,
			Body: body,
		}
		_article.Create()
		fmt.Println(_article)
		if  _article.ID > 0 {
			fmt.Println("出插入成功，ID为"+strconv.FormatInt(_article.ID,10))
			showUrl := route.Name2URL("articles.index")
			http.Redirect(w,r,showUrl,http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"创建文章失败，请联系管理员")
		}
	} else {
		storeURL := route.Name2URL("articles.store")
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
