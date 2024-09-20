package requests

import "time"

// ArticleDetailRes
// @Description: 文章详情响应架构体
type ArticleDetailRes struct {
	Id              uint       `json:"id"`
	Title           string     `json:"title"`
	LikeCount       int        `json:"like_count"`
	CollectionCount int        `json:"collection_count"`
	ViewsCount      int        `json:"views_count"`  // 浏览量
	Heat            int        `json:"heat"`         // 文章热度
	CategoryID      uint       `json:"category_id"`  // 所属类目ID，外键
	Summary         string     `json:"summary"`      // 文章摘要
	PublishedAt     *time.Time `json:"published_at"` // 发布时间 可以为空(如草稿)
	Content         string     `json:"content"`      // 文章内容

	ImageUrl string `json:"image_url"` // 文章封面url
	Author   string `json:"author"`
}