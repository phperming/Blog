package article

import (
	"Blog/pkg/model"
	"Blog/pkg/route"
	"Blog/pkg/types"
	"strconv"
)

type Article struct {
	ID int64
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

func (article Article)Create() (err error) {
	if err := model.DB.Create(&article).Error; err != nil {
		return err
	}

	return nil
}

func (a Article)Link() string {
	return route.Name2URL("articles.show","id",strconv.FormatInt(a.ID,10))
}
