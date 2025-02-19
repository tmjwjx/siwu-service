package models

import "gorm.io/gorm"

type CommentLike struct {
	gorm.Model      //ID CreatedAt UpdatedAt DeletedAt
	CommentID  uint `json:"comment_id" gorm:"index"`      // 评论ID
	UserID     uint `json:"user_id" gorm:"index"`         // 点赞人ID
	IsRead     bool `json:"is_read" gorm:"default:false"` // 是否已读
}
