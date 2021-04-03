package view

import (
	"Blog/pkg/logger"
	"Blog/pkg/route"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

//Render 渲染视图
func Render(w io.Writer,name string,data interface{})  {
	//1.设置模板的相对路径
	viewDir := "resource/views/"
	//2.语法糖，将 article.show 更正为 articles/show
	name = strings.Replace(name, ".", "/", -1)
	//3.所有布局文件的Slice
	files,err := filepath.Glob(viewDir+"/layouts/*.gohtml")
	logger.LogError(err)
	//4.在Slice里面增加我们的目标文件
	newFiles := append(files,viewDir+name+".gohtml")
	//5.解析所有的模板文件
	tmpl, err := template.New(name + ".gohtml").Funcs(template.FuncMap{
		"RouteName2URL": route.Name2URL,
	}).ParseFiles(newFiles...)
	logger.LogError(err)
	//6.渲染模板
	tmpl.ExecuteTemplate(w,"app",data)
}
