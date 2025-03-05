package models

import "gorm.io/gorm"

type ArticleView struct {
	gorm.Model
	ArticleID uint `json:"article_id" gorm:"index"` // 文章ID
	UserID    uint `json:"user_id" gorm:"index"`    // 用户ID
}
