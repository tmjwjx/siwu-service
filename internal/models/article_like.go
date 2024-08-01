package models

import "gorm.io/gorm"

// ArticleLike 文章获赞表
type ArticleLike struct {
	gorm.Model      //ID CreatedAt UpdatedAt DeletedAt
	ArticleID  uint `json:"article_id" gorm:"index"` // 文章ID
	UserID     uint `json:"user_id" gorm:"index"`    // 点赞人ID
}
