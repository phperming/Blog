package article

import (
	"Blog/pkg/model"
	"Blog/pkg/types"
)

type Article struct {
	ID int
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
