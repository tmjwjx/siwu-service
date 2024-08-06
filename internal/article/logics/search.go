package logics

import (
	"forum/internal/article/repositories"
	"forum/internal/article/requests"
	"forum/internal/models"
	"gorm.io/gorm"
)

func Search(db *gorm.DB, req requests.ReqSearch) ([]models.Article, error) {
	
	articles, err := repositories.SearchArticles(db, req)
	if err != nil {
		return nil, err
	}
	return articles, nil
}
func Public(db *gorm.DB, req requests.ReqPublish) error {
	
	return nil
}