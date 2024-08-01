package models

import (
	"gorm.io/gorm"
	"time"
)

// Article 文章表
type Article struct {
	gorm.Model                             //ID CreatedAt UpdatedAt DeletedAt
	UserID             uint                `json:"user_id" gorm:"index"` // 作者ID，外键
	Title              string              `json:"title"`                // 文章标题
	Image              string              `json:"image"`                // 文章图片
	Link               string              `json:"link"`                 // 文章链接，唯一 (链接后缀为：作者id+文章id+类目+发布日期)
	LikesCount         int                 `json:"likes_count"`          // 点赞数量
	CollectionsCount   int                 `json:"collections_count"`    // 收藏数量
	CommentsCount      int                 `json:"comments_count"`       // 评论数量
	ViewsCount         int                 `json:"views_count"`          // 浏览量
	Heat               int                 `json:"heat"`                 // 文章热度
	Status             string              `json:"status"`               // 文章属性（草稿，私有，公开）
	CategoryID         uint                `json:"category_id"`          // 所属类目ID，外键
	Summary            string              `json:"summary"`              // 文章摘要
	PublishedAt        *time.Time          `json:"published_at"`         // 发布时间
	Tags               []Tag               `gorm:"many2many:article_tags"`
	ArticleLikes       []ArticleLike       `gorm:"many2many:user_article_likes"`      //点赞的人
	ArticleCollections []ArticleCollection `gorm:"many2many:user_article_collection"` // 收藏的文章
	ArticleComments    []ArticleComment    `gorm:"many2many:article_comments"`        // 文章的评论
}
