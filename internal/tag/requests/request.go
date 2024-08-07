package requests

// TagReq 请求标签结构体
type TagReq struct {
	ID uint `json:"id" form:"id"`
}

// TagRes 响应标签结构体
type TagRes struct {
	ID           uint     `json:"id" form:"id"`
	Name         string   `json:"name" form:"name"`                   // 标签名称
	Description  string   `json:"description" form:"description"`     // 标签描述
	ArticleCount string   `json:"article_count" form:"article_count"` // 标签关联的文章数量
	Heat         string   `json:"heat" form:"heat"`                   // 标签热度
	FansCount    string   `json:"fans_count" form:"fan_count"`        // 关注人数
	Path         []string `json:"path" form:"path"`
}
