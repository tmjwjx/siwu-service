package logics

import (
	"forum/internal/article/repositories"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

// ArticleDetailLogic
// @Description: 文章详情
// @param        db *gorm.DB
// @param        id string
// @return       data
// @return       err
func ArticleDetailLogic(db *gorm.DB, id string) (data interface{}, err error) {

	data, err = repositories.ArticleDetailRep(db, id)
	if err != nil {
		globals.Log.Errorf("err = %s", err)
		return nil, err
	}
	return nil, nil

}