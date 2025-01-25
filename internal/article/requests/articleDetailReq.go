package requests

import (
	"forum/internal/models"
	"time"
)

// ArticleDetailRes
// @Description: 文章详情响应结构体
type ArticleDetailRes struct {
	Id               uint         `json:"id"`                                                           // 文章ID
	Title            string       `json:"title"`                                                        // 文章标题
	UserId           uint         `json:"user_id"`                                                      // 作者ID
	LikesCount       int          `json:"likes_count"`                                                  // 点赞数
	CollectionsCount int          `json:"collections_count"`                                            // 收藏数
	ViewsCount       int          `json:"views_count"`                                                  // 浏览量
	Heat             int          `json:"heat"`                                                         // 文章热度
	LikeStatus       bool         `json:"like_status"`                                                  // 点赞状态
	CollectionStatus bool         `json:"collection_status"`                                            // 收藏状态
	CategoryID       uint         `json:"category_id"`                                                  // 所属类目ID，外键
	Summary          string       `json:"summary"`                                                      // 文章摘要
	PublishedAt      *time.Time   `json:"published_at"`                                                 // 发布时间 可以为空(如草稿)
	FormatTime       string       `json:"format_time"`                                                  // 格式化时间
	DailyTime        string       `json:"daily_time"`                                                   // 日常时间
	Content          string       `json:"content"`                                                      // 文章内容
	ImageUrl         string       `json:"image_url"`                                                    // 文章封面url
	Nickname         string       `json:"nickname"`                                                     // 作者昵称
	Tags             []models.Tag `json:"tags" gorm:"many2many:article_tags;joinForeignKey:article_id"` // 标签
}

// AboutArticleRes
// @Description: 相关文章 响应结构体
type AboutArticleRes struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	ViewsCount int    `json:"views_count"`
	LikesCount int    `json:"likes_count"`
}

// ArticleLikeReq
// @Description: 文章点赞请求结构体
type ArticleLikeReq struct {
	ArticleId  uint `json:"article_id" form:"article_id"`
	LikeStatus bool `json:"like_status" form:"like_status"`
}

// ArticleCollectionReq
// @Description: 文章收藏请求结构体
type ArticleCollectionReq struct {
	ArticleId        uint `json:"article_id" form:"article_id"`
	CollectionStatus bool `json:"collection_status" form:"collection_status"`
}
