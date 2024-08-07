package logics

import (
	"forum/internal/article/repositories"
	"forum/internal/article/requests"
	"forum/internal/models"
	"gorm.io/gorm"
)

// ArticleSearchLogic 搜索文章
func ArticleSearchLogic(db *gorm.DB, req requests.ReqSearch) ([]models.Article, error) {

	articles, err := repositories.QueryArticlesRep(db, req)
	if err != nil {
		return nil, err
	}
	return articles, nil
}