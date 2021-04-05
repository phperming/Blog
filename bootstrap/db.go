package bootstrap

import (
	"Blog/app/models/article"
	"Blog/app/models/user"
	"Blog/pkg/model"
	"gorm.io/gorm"
	"time"
)

//初始化数据库和ORM
func SetupDB() {
	//建立数据库连接池
	db := model.ConnectDB()

	//命令行打印数据库请求的的信息
	sqlDB, _ := db.DB()

	//设置最大连接数
	sqlDB.SetMaxOpenConns(100)
	//设置最大空闲连接数
	sqlDB.SetMaxIdleConns(25)
	//设置每个连接的过期时间
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	migration(db)
	
}

func migration(db *gorm.DB)  {
	//自动迁移
	db.AutoMigrate(
		&user.User{},
		&article.Article{},
	)
}
