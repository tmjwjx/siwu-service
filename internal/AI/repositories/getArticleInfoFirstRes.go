package repositories

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"gorm.io/gorm"
)

func UpdateArticleAbstract(db *gorm.DB, abstract string, articleID uint) error {
	var article models.Article
	err := db.Model(&models.Article{}).Select("summary").Where("id = ?", articleID).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("数据库中没有查询到该文章")
		} else {
			return fmt.Errorf("UpdateArticleAbstract -> %v", err)
		}
	}

	if article.Summary == "" {
		result := db.Model(&models.Article{}).Where("id = ?", articleID).Update("summary", abstract)
		if result.Error != nil {
			return fmt.Errorf("UpdateArticleAbstract -> %v", err)
		} else if result.RowsAffected > 0 {
			return nil
		}
	}

	return nil
}
