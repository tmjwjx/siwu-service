package requests

import (
	"forum/internal/models"
	"time"
)

// ArticleListReq
// @Description: 获取文章列表请求
type ArticleListReq struct {
	Page             int       `json:"page"`              // 分页页码
	Limit            int       `json:"limit"`             // 每页条数
	ArticleCondition int       `json:"article_condition"` // 文章状态（是否封禁）
	StartTime        time.Time `json:"startTime"`         // 发布开始时间
	EndTime          time.Time `json:"endTime"`           // 发布截至时间
	Keyword          string    `json:"keyword"`           // 文章标题
	ArticleTags      []int     `json:"article_tags"`      // 文章标签id，数组类型
	Nickname         string    `json:"nickname"`          // 发布人用户名
	ViewsCount       int       `json:"views_count"`       // 文章浏览量
	LikesCount       int       `json:"likes_count"`       // 文章点赞量
	CollectionsCount int       `json:"collections_count"` // 文章收藏量
	CommentsCount    int       `json:"comments_count"`    // 文章评论量
	Kind             int       `json:"kind"`              // 排序方式
	Heat             int       `json:"heat"`              // 文章热度
}

// SearchArticleListRes
// @Description: 获取文章列表响应
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
	Summary          string `json:"summary"`           // 文章摘要
	ImageUrl         string `json:"image_url"`         // 文章封面url
	PublishedAt      string `json:"published_at"`      // 发布时间
	UpdatedAt        string `json:"updated_at"`        // 更新时间
	Status           string `json:"status"`            // 文章属性（草稿，私有，公开）
	//Tags             []Tags `json:"tags" gorm:"many2many:article_tags;joinForeignKey:article_id;table:sw_tags"`
	Tags []models.Tag `json:"tags" gorm:"many2many:article_tags;joinForeignKey:article_id"`
}

//type Tags struct {
//	ID   int    `json:"id"`   // 标签ID
//	Name string `json:"name"` // 标签名称
//}

// GetArticleByTagReq
// @Description: 获取标签下的文章请求
// @Author tianjiajie 2025-01-16 14:38:01
type GetArticleByTagReq struct {
	Id    int `json:"id" form:"id"`       // 标签ID
	Kind  int `json:"kind" form:"kind"`   // 类型
	Page  int `json:"page" form:"page"`   // 分页页码
	Limit int `json:"limit" form:"limit"` // 每页条数
}

type ArticleOperationListReq struct {
	IdList []int `json:"id_list"`
}
