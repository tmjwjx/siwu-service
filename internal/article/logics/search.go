package logics

import (
	"forum/internal/article/repositories"
	"forum/internal/article/requests"
	"forum/internal/models"
	"gorm.io/gorm"
)

func Search(db *gorm.DB, req requests.SearchRequest) ([]models.Article, error) {
	
	offset := (req.Page - 1) * req.Limit // 计算当前页的偏移量，用于分页
	
	articles, err := repositories.SearchArticles(db, req, offset)
	if err != nil {
		return nil, err
	}
	return articles, nil
}