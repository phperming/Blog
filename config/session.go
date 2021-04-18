package config

import "Blog/pkg/config"

func init()  {
	config.Add("session", config.StrMap{
		//目前支支持Cookie
		"default" : config.Env("SESSION_DRIVER","cookie"),

		//会话的COOKIE名称
		"session_name" : config.Env("SESSION_NAME","goblog-session"),
	})
}