package models

import "time"

// Article 文章表
type Article struct {
	ID                      uint       `json:"id" gorm:"primaryKey"`                                                                                                   // 主键
	AuthorID                uint       `json:"author_id" gorm:"index"`                                                                                                 // 作者ID，外键
	ArticleTitle            string     `json:"article_title"`                                                                                                          // 文章标题
	ArticleImage            string     `json:"article_image"`                                                                                                          // 文章图片
	ArticleLink             string     `json:"article_link"`                                                                                                           // 文章链接，唯一 (链接后缀为：作者id+文章id+类目+发布日期)
	ArticleLikesCount       int        `json:"article_likes_count"`                                                                                                    // 点赞数量
	ArticleCollectionsCount int        `json:"article_collections_count"`                                                                                              // 收藏数量
	ArticleCommentsCount    int        `json:"article_comments_count"`                                                                                                 // 评论数量
	ArticleViewsCount       int        `json:"article_views_count"`                                                                                                    // 浏览量
	ArticleHeat             int        `json:"article_heat"`                                                                                                           // 文章热度
	ArticleAttributes       string     `json:"article_attributes"`                                                                                                     // 文章属性（草稿，私有，公开）
	ArticleCategoryID       uint       `json:"article_category_id"`                                                                                                    // 所属类目ID，外键
	ArticleTagID            []Tag      `json:"article_tag_id" gorm:"many2many:article_tags;foreignKey:ID;joinForeignKey:ArticleID;references:ID;joinReferences:TagID"` // 文章标签ID
	ArticleSummary          string     `json:"article_summary"`                                                                                                        // 文章摘要
	ArticleContent          string     `json:"article_content"`                                                                                                        // 文章内容
	UpdatedAt               *time.Time `json:"updated_at"`                                                                                                             // 编辑时间
	PublishedAt             *time.Time `json:"published_at"`                                                                                                           // 发布时间
	Category                Category   `gorm:"foreignKey:ArticleCategoryID"`
}

// Category 类目表
type Category struct {
	ID            uint   `json:"id" gorm:"primaryKey"` // 主键
	CategoryName  string `json:"category_name"`        // 类目名称
	CategoryImage string `json:"category_image"`       // 类目图标
}

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

// ArticleLike 文章点赞表
type ArticleLike struct {
	ID        uint    `json:"id" gorm:"primaryKey"`    // 主键
	ArticleID uint    `json:"article_id" gorm:"index"` // 文章ID
	UserID    uint    `json:"user_id" gorm:"index"`    // 点赞人ID
	Article   Article `gorm:"foreignKey:ArticleID"`
}

// ArticleCollection 文章收藏表
type ArticleCollection struct {
	ID        uint    `json:"id" gorm:"primaryKey"`    // 主键
	ArticleID uint    `json:"article_id" gorm:"index"` // 文章ID
	UserID    uint    `json:"user_id" gorm:"index"`    // 收藏人ID
	Article   Article `gorm:"foreignKey:ArticleID"`
}

// Attachment 图片附件表
type Attachment struct {
	Id          uint    `json:"id" gorm:"primaryKey"` // 主键
	FileName    string  `json:"file_name"`            // 文件名
	FileType    string  `json:"file_type"`            // 文件类型，例如 image/png
	FileSize    int64   `json:"file_size"`            // 文件大小（以字节为单位）
	FileContent []byte  `json:"file_content"`
	ArticleID   uint    `json:"article_id"` // 所属文章 ID，外键
	Article     Article `gorm:"foreignKey:ArticleID"`
}
