package models

import "gorm.io/gorm"

// ArticleCollection 文章收藏表
type ArticleCollection struct {
	gorm.Model      //ID CreatedAt UpdatedAt DeletedAt
	ArticleID  uint `json:"article_id" gorm:"index"` // 文章ID
	UserID     uint `json:"user_id" gorm:"index"`    // 收藏人ID
}