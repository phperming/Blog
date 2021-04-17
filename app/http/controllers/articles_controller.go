package controllers

import (
	"Blog/app/models/article"
	"Blog/pkg/logger"
	"Blog/pkg/route"
	"Blog/pkg/view"
	"fmt"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"unicode/utf8"
)

type ArticlesController struct {

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
		fmt.Println("chenggong")
		view.Render(w,view.D{
			"Article" : article,
		},"articles.show")
	}
}

func (*ArticlesController)Index(w http.ResponseWriter,r *http.Request)  {
	articles,err := article.GetAll()
	//fmt.Println(articles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w,"500 服务器内部错误")
	} else {
		view.Render(w,view.D{
			"Articles" : articles,
		},"articles.index")
	}
}

func (*ArticlesController)Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w,view.D{},"articles.create","articles._form_field")
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
			fmt.Println("出插入成功，ID为"+_article.GetStringID())
			showUrl := route.Name2URL("articles.index")
			http.Redirect(w,r,showUrl,http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"创建文章失败，请联系管理员")
		}
	} else {
		view.Render(w,view.D{
			"Title": title,
			"Body": body,
			"Errors": errors,
		},"articles.create","articles._form_field")
	}

}

func (*ArticlesController)Edit(w http.ResponseWriter,r *http.Request) {
	//获取URL参数
	id := route.GetRouterVariable("id",r)

	//读取对应的文章数据
	_article,err := article.Get(id)

	//如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound{
			//数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w,"404文章未找到")
		} else {
			//数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"500 服务器内部错误")
		}
	} else {
		//读取成功，显示表单
		fmt.Println("读取成功")
		view.Render(w,view.D{
			"Title": _article.Title,
			"Body": _article.Body,
			"Article" : _article,
			"Errors": nil,
		},"articles.edit","articles._form_field")
	}

}

func (*ArticlesController)Update(w http.ResponseWriter,r *http.Request) {
	//获取文章ID
	id := route.GetRouterVariable("id",r)

	//获取文章
	_article,err := article.Get(id)

	//如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w,"404 文章未找到")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"500 服务器内部错误")
		}
	} else {
		//未出现错误
		//验证表单
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errors := validateArticleFormData(title,body)

		if len(errors) == 0 {
			//验证通过
			_article.Title = title
			_article.Body = body

			rowsAffected,err := _article.Update()

			if err != nil {
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w,"500 服务器内部错误")
			}

			//更新成功，跳转到文章详情页
			if rowsAffected > 0 {
				showUrl := route.Name2URL("articles.show","id",id)
				http.Redirect(w,r,showUrl,http.StatusFound)
			} else  {
				fmt.Fprint(w,"没有做任何更改")
			}
		} else {
			//表单验证不通过显示理由

			data := view.D{
				"Title": title,
				"Body": body,
				"Article": _article,
				"Errors": errors,
			}

			tmpl ,err := template.ParseFiles("resource/views/articles/edit.gohtml")

			logger.LogError(err)
			tmpl.Execute(w,data)
		}

	}
}

func (*ArticlesController)Delete(w http.ResponseWriter,r *http.Request) {
	//获取id
	id := route.GetRouterVariable("id", r)

	//读取对应的文章
	_article, err := article.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w,"404 文章未找到")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"500 服务器内部错误")
		}
	} else {
		//未出现错误，执行删除操作
		rowsAffected,err := _article.Delete()

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
				indexURL := route.Name2URL("articles.index")
				http.Redirect(w,r,indexURL,http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w,"404 文章未找到")
			}
		}
	}
}

func validateArticleFormData(title string,body string) map[string]string {
	errors := make(map[string]string)

	//验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度不能小于10"
	}

	//验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需介于3-40之间"
	}

	return errors
}
