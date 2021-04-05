package requests

import (
	"Blog/app/models/user"
	"github.com/thedevsaddam/govalidator"
)

func ValidateRegistrationForm(data user.User) map[string][]string  {
	//1.定制规则
	rules := govalidator.MapData{
		"name" : []string{"required","alpha_num","between:3,20"},
		"email" : []string{"required","min:4","max:30","email"},
		"password" : []string{"required","min:4"},
		"password_confirm" : []string{"required"},
	}

	//2.定制错误信息
	messages := govalidator.MapData{
		"name" : []string{
			"required:用户名为必填项",
			"alpha_num:格式错误，只允许数字",
			"between:用户名长度必须在3-30之间",
		},
		"email": []string{
			"required:邮箱为必填项",
			"min:Email 长度需大于 4",
			"max:Email 长度需小于 30",
			"email:Email 格式不正确，请提供有效的邮箱地址",
		},
		"password": []string{
			"required:密码为必填项",
			"min:密码长度需大于6",
		},
		"password_confirm": []string{
			"required:确认密码为必填项",
		},
	}

	//3.配置初始化
	opts := govalidator.Options{
		Data: &data,
		Rules: rules,
		TagIdentifier: "valid",
		Messages: messages,
	}

	//4.开始认证、
	errs := govalidator.New(opts).ValidateStruct()

	// 5. 因 govalidator 不支持 password_confirm 验证，我们自己写一个
	if data.Password != data.PasswordConfirm {
		errs["password_confirm"] = append(errs["password_confirm"],"两次输入的密码不匹配")
	}

	return errs
}
