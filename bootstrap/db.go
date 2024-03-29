package bootstrap

import (
	"Blog/app/models/article"
	"Blog/app/models/user"
	"Blog/pkg/config"
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
	sqlDB.SetMaxOpenConns(config.GetInt("database.mysql.max_open_connections"))
	//设置最大空闲连接数
	sqlDB.SetMaxIdleConns(config.GetInt("database.mysql.max_idle_connections"))
	//设置每个连接的过期时间
	sqlDB.SetConnMaxLifetime(time.Duration(config.GetInt("database.mysql_max_life_seconds")) * time.Second)

	migration(db)
	
}

func migration(db *gorm.DB)  {
	//自动迁移
	db.AutoMigrate(
		&user.User{},
		&article.Article{},
	)
}
