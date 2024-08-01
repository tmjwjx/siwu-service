package models

import "time"

// Comment 评论表
type Comment struct {
	ID                 uint       `json:"id" gorm:"primaryKey"`           // 主键
	ArticleID          uint       `json:"article_id" gorm:"index"`        // 所属文章ID，外键
	UserID             uint       `json:"user_id" gorm:"index"`           // 作者ID
	CommentParentID    uint       `json:"comment_parent_id" gorm:"index"` // 上一条评论ID
	CommentPublishedAt *time.Time `json:"comment_published_at"`           // 发布时间
	CommentContent     string     `json:"comment_content"`                // 评论内容
	CommentLikesCount  int        `json:"comment_likes_count"`            // 点赞数量
	Article            Article    `gorm:"foreignKey:ArticleID"`
}
