package logics

import (
	"forum/internal/article/repositories"
	"forum/internal/article/requests"
	"forum/internal/models"
	"gorm.io/gorm"
)

func ArticleSearch(db *gorm.DB, req requests.ReqSearch) ([]models.Article, error) {

	articles, err := repositories.QueryArticles(db, req)
	if err != nil {
		return nil, err
	}
	return articles, nil
}