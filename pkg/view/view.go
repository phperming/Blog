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
func Render(w io.Writer,data interface{},tplFiles...string)  {
	//1.设置模板的相对路径
	viewDir := "resource/views/"
	//2.语法糖，将 article.show 更正为 articles/show
	for i,f := range tplFiles {
		tplFiles[i] = viewDir + strings.Replace(f,".","/",-1) + ".gohtml"
	}
	//3.所有布局文件的Slice
	layoutFiles,err := filepath.Glob(viewDir+"/layouts/*.gohtml")
	logger.LogError(err)
	//4.在Slice里面增加我们的目标文件
	allFiles := append(layoutFiles,tplFiles...)
	//5.解析所有的模板文件
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"RouteName2URL": route.Name2URL,
	}).ParseFiles(allFiles...)
	logger.LogError(err)
	//6.渲染模板
	tmpl.ExecuteTemplate(w,"app",data)
}
