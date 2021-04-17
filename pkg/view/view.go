package view

import (
	"Blog/pkg/auth"
	"Blog/pkg/flash"
	"Blog/pkg/logger"
	"Blog/pkg/route"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

type D map[string]interface{}

//Render 渲染视图
func Render(w io.Writer,data D,tplFiles...string)  {
	RenderTemplate(w,"app",data,tplFiles...)
}

func RendSimple(w io.Writer,data D,tplFiles ...string)  {
	RenderTemplate(w,"simple",data,tplFiles...)
}

//渲染视图
func RenderTemplate(w io.Writer,name string,data D,tplFiles ...string)  {
	//1.通用模板数据
	data["isLogined"] = auth.Check()
	data["loginUser"] = auth.User
	data["flash"] = flash.All()

	//2.生成模板文件
	allFiles := getTemplateFiles(tplFiles...)

	//3.解析所有的模板文件
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"RouteName2URL": route.Name2URL,
	}).ParseFiles(allFiles...)
	logger.LogError(err)

	//4.渲染模板
	tmpl.ExecuteTemplate(w,name,data)
}

func getTemplateFiles(tplFiles ...string) []string  {
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

	return allFiles
}
