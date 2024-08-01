package models

import (
	"gorm.io/gorm"
)

// ArticleComment 评论表
type ArticleComment struct {
	gorm.Model        //ID CreatedAt UpdatedAt DeletedAt
	ArticleID  uint   `json:"article_id" gorm:"index"` // 所属文章ID，外键
	UserID     uint   `json:"user_id" gorm:"index"`    // 作者ID
	HighestID  uint   `json:"highest_id"`              // 最上层一级评论
	ParentID   uint   `json:"parent_id" gorm:"index"`  // 上一条评论ID
	Content    string `json:"content"`                 // 评论内容
	LikesCount int    `json:"likes_count"`             // 点赞数量
	// 评论点赞还需要一个表
	Users []User `gorm:"many2many:comment_likes"` // 用户对评论的点赞
}
