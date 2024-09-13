package requests

import (
	"forum/internal/models"
	"time"
)

// ArticleListReq
// @Description: 获取文章列表
type ArticleListReq struct {
	Page             int       `json:"page"`              // 分页页码
	Limit            int       `json:"limit"`             // 每页条数
	ArticleCondition int       `json:"article_condition"` // 文章状态（是否封禁）
	StartTime        time.Time `json:"startTime"`         // 发布开始时间
	EndTime          time.Time `json:"endTime"`           // 发布截至时间
	Title            string    `json:"title"`             // 文章标题
	ArticleTags      []int     `json:"article_tags"`      // 文章标签id，数组类型
	Nickname         string    `json:"nickname"`          // 发布人用户名
	ViewsCount       int       `json:"views_count"`       // 文章浏览量
	LikesCount       int       `json:"likes_count"`       // 文章点赞量
	CollectionsCount int       `json:"collections_count"` // 文章收藏量
	CommentsCount    int       `json:"comments_count"`    // 文章评论量
	Heat             int       `json:"heat"`              // 文章热度
}

type SearchArticleListRes struct {
	ID               uint   `json:"id"`                // 文章ID
	Title            string `json:"title"`             // 文章标题
	ArticleCondition int    `json:"article_condition"` // 文章状态
	ViewsCount       int    `json:"views_count"`       // 浏览数
	LikesCount       int    `json:"likes_count"`       // 点赞数
	CollectionsCount int    `json:"collections_count"` // 收藏数
	CommentsCount    int    `json:"comments_count"`    // 评论数
	Heat             int    `json:"heat"`              // 热度
	Nickname         string `json:"nickname"`          // 发布人用户名
	PublishedAt      string `json:"published_at"`      // 发布时间
	UpdatedAt        string `json:"updated_at"`        // 更新时间
	//Tags             []Tags `json:"tags" gorm:"many2many:article_tags;joinForeignKey:article_id;table:sw_tags"`
	Tags []models.Tag `json:"tags" gorm:"many2many:article_tags;joinForeignKey:article_id"`
}

type Tags struct {
	ID   int    `json:"id"`   // 标签ID
	Name string `json:"name"` // 标签名称
}

//type ArcList struct {
//	gorm.Model //ID CreatedAt UpdatedAt DeletedAt
//	//UserID           uint         `json:"user_id" gorm:"index"`                                                // 作者ID，外键
//	//Title            string       `json:"title" gorm:"not null"`                                               // 文章标题
//	//CollectionsCount int          `json:"collections_count" gorm:"default:0"`                                  // 收藏数量
//	//CommentsCount    int          `json:"comments_count" gorm:"default:0"`                                     // 评论数量
//	//ViewsCount       int          `json:"views_count" gorm:"default:0"`                                        // 浏览量
//	//Heat             int          `json:"heat" gorm:"default:0"`                                               // 文章热度
//	//Status           string       `json:"status" gorm:"type:enum('draft','private','public');default:'draft'"` // 文章属性（草稿，私有，公开）
//	//CategoryID       uint         `json:"category_id" gorm:"default:0;index"`                                  // 所属类目ID，外键
//	//Summary          string       `json:"summary" gorm:"type:text"`                                            // 文章摘要
//	//PublishedAt      *time.Time   `json:"published_at"`                                                        // 发布时间 可以为空(如草稿)
//	//Content          string       `json:"content" gorm:"type:text;not null"`                                   // 文章内容
//	//ArticleCondition int          `json:"article_condition"`                                                   // 是否封禁 1：正常 2：封禁
//	Tags []models.Tag `gorm:"many2many:article_tags;joinForeignKey:article_id"`
//}