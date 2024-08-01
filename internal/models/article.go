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
