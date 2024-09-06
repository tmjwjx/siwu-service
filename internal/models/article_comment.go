package models

import (
	"gorm.io/gorm"
)

// ArticleComment 评论表
type ArticleComment struct {
	gorm.Model          //ID CreatedAt UpdatedAt DeletedAt
	ArticleID    uint   `json:"article_id" gorm:"index"`      // 所属文章ID，外键
	UserID       uint   `json:"user_id" gorm:"index"`         // 作者ID
	HighestID    *uint  `json:"highest_id" gorm:"index"`      // 最上层一级评论
	ParentID     *uint  `json:"parent_id" gorm:"index"`       // 上一条评论ID , 允许为 null，表示顶级评论
	ParentUserID *uint  `json:"parent_user_id"`               // 上一条评论的发布用户ID
	Content      string `json:"content" gorm:"not null"`      // 评论内容
	LikesCount   int    `json:"likes_count" gorm:"default:0"` // 点赞数量
}