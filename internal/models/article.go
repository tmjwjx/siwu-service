package models

import (
	"gorm.io/gorm"
	"time"
)

// Article 文章表
type Article struct {
	gorm.Model                  //ID CreatedAt UpdatedAt DeletedAt
	UserID           uint       `json:"user_id" gorm:"index"`                                                // 作者ID，外键
	Title            string     `json:"title" gorm:"not null"`                                               // 文章标题
	Link             string     `json:"link" gorm:"not null"`                                                // 文章链接，唯一 (链接后缀为：作者id+文章id+类目+发布日期)
	LikesCount       int        `json:"likes_count" gorm:"default:0"`                                        // 点赞数量
	CollectionsCount int        `json:"collections_count" gorm:"default:0"`                                  // 收藏数量
	CommentsCount    int        `json:"comments_count" gorm:"default:0"`                                     // 评论数量
	ViewsCount       int        `json:"views_count" gorm:"default:0"`                                        // 浏览量
	Heat             int        `json:"heat" gorm:"default:0"`                                               // 文章热度
	Status           string     `json:"status" gorm:"type:enum('draft','private','public');default:'draft'"` // 文章属性（草稿，私有，公开）
	CategoryID       uint       `json:"category_id" gorm:"default:0;index"`                                  // 所属类目ID，外键
	Summary          string     `json:"summary" gorm:"type:text"`                                            // 文章摘要
	PublishedAt      *time.Time `json:"published_at"`                                                        // 发布时间 可以为空(如草稿)
	Content          string     `json:"content" gorm:"type:text;not null"`                                   // 文章内容
	ArticleCondition int        `json:"article_condition"`                                                   // 是否封禁 1：正常 2：封禁
	Tags             []Tag      `gorm:"many2many:article_tags"`
	UserLikes        []User     `gorm:"many2many:article_likes"`       //点赞的人
	UserCollections  []User     `gorm:"many2many:article_collections"` // 收藏的文章
	ArticleComments  []User     `gorm:"many2many:article_comments"`    // 文章的评论
}