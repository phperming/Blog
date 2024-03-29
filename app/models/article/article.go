package article

import (
	"Blog/app/models"
	"Blog/pkg/logger"
	"Blog/pkg/model"
	"Blog/pkg/route"
	"Blog/pkg/types"
)

type Article struct {
	models.BaseModel
	Title string
	Body string
}


func Get(idstr string)(Article,error)  {
	var article Article
	id := types.StringToInt(idstr)

	if err := model.DB.First(&article, id).Error;err != nil {
		return article,err
	}

	return article,nil
}

func (article *Article)Create() (err error) {
	if err := model.DB.Create(&article).Error; err != nil {
		return err
	}

	return nil
}

func (article *Article)Update() (rowsAffected int64,err error) {
	result := model.DB.Save(&article)
	if err := result.Error; err != nil {
		logger.LogError(err)
		return 0,err
	}

	return result.RowsAffected,nil
}

func (article *Article)Delete() (rowsAffected int64,err error) {
	result := model.DB.Delete(&article)
	if err := result.Error; err != nil {
		logger.LogError(err)
		return 0,err
	}

	return result.RowsAffected,nil
}

func (a Article)Link() string {
	return route.Name2URL("articles.show","id",a.GetStringID())
}
