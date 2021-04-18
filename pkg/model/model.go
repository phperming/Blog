package model

import (
	"Blog/pkg/config"
	"Blog/pkg/logger"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB{
	var err error

	//初始化数据库连接信息
	var(
		host = config.Env("database.mysql.host")
		port = config.Env("database.mysql.port")
		database = config.Env("database.mysql.database")
		username = config.Env("database.mysql.username")
		password = config.Env("database.mysql.password")
		charset = config.Env("database.mysql.charset")
	)

	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		username,password,host,port,database,charset,true,"Local")
	//config := mysql.New(mysql.Config{
	//	DSN: "root:root@tcp(127.0.0.1:3306)/goblog?charset=utf8&parseTime=True&loc=Local",
	//})
	gormCofig := mysql.New(mysql.Config{
		DSN: dns,
	})

	var level gormlogger.LogLevel
	if config.GetBool("app.debug") {
		//读取不到数据也会显示
		level = gormlogger.Warn
	} else {
		//只有错误才会显示
		level = gormlogger.Error
	}

	//准备数据库连接池
	DB, err = gorm.Open(gormCofig, &gorm.Config{
		Logger: gormlogger.Default.LogMode(level),
	})
	logger.LogError(err)

	return DB
}
