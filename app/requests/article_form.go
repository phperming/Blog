package requests

import (
	"Blog/app/models/article"
	"github.com/thedevsaddam/govalidator"
)

func ValidateArticleForm(data article.Article) map[string][]string  {
	//1.定制认证规则
	rules := govalidator.MapData{
		"title" : []string{"required","min:3","max:40"},
		"body" : []string{"required","min:10"},
	}

	//定制错误信息
	messages := govalidator.MapData{
		"title" : []string{
			"required:标题为必填项",
			"min:标题长度需大于3个字符",
			"max:标题长度需小于40",
		},
		"body":[]string{
			"required:内容为必填项",
			"min:长度需大于10",
		},
	}

	//配置初始化
	opts := govalidator.Options{
		Data : &data,
		Rules: rules,
		TagIdentifier: "valid",
		Messages: messages,
	}

	//开始验证
	return govalidator.New(opts).ValidateStruct()
}
