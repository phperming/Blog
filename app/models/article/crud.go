package article

import "Blog/pkg/model"

func GetAll() ([]Article,error) {
	var articles []Article

	if err := model.DB.Find(&articles).Error;err != nil {
		return articles, err
	}

	return articles,nil
}
