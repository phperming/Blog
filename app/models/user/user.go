package user

import (
	"Blog/app/models"
	"Blog/pkg/model"
	"Blog/pkg/types"
)

type User struct {
	models.BaseModel

	Name string `gorm:"type:varchar(50);not null;unique" valid:"name"`
	Email string `gorm:"type:varchar(100);unique" valid:"email"`
	Password string `gorm:"type:varchar(100)" valid:"password"`
	// gorm:"-" —— 设置 GORM 在读写时略过此字段
	PasswordConfirm string `gorm:"-" valid:"password_confirm"`
}

//对比密码是否匹配
func (u User)ComparePassword(password string) bool  {
	return u.Password == password
}

func Get(uid string)  (User,error) {
	var user User
	_uid := types.StringToInt(uid)
	if err := model.DB.First(&user,_uid).Error; err != nil {
		return user, err
	}

	return user,nil
}

func GetByEmail(email string) (User,error) {
	var user User
	if err := model.DB.Where("email=?",email).First(&user).Error; err != nil {
		return user,err
	}

	return user,nil
}







