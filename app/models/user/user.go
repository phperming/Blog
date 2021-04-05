package user

import (
	"Blog/app/models"
)

type User struct {
	models.BaseModel

	Name string `gorm:"type:varchar(50);not null;unique" valid:"name"`
	Email string `gorm:"type:varchar(100);unique" valid:"email"`
	Password string `gorm:"type:varchar(100)" valid:"password"`
	// gorm:"-" —— 设置 GORM 在读写时略过此字段
	PasswordConfirm string `gorm:"-" valid:"password_confirm"`
}





