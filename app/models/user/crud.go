package user

import (
	"Blog/pkg/logger"
	"Blog/pkg/model"
)

func (user *User)Create() (error) {
	if err := model.DB.Create(&user).Error;err != nil {
		logger.LogError(err)
		return err
	}

	return nil

}
