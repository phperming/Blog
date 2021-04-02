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
		tmpl,err := template.New("show.gohtml").Funcs(template.FuncMap{
			"RouteName2URL" : route.Name2URL,
			"Int64ToString" : types.Int64ToString,
		}).ParseFiles("resource/views/articles/show.gohtml")
		//tmpl, err := template.ParseFiles("resource/views/articles/show.gohtml")
		logger.LogError(err)
		tmpl.Execute(w,article)
	}
}
