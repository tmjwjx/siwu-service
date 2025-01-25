package requests

// ReqPublish
// @Description: 发布文章请求
type ReqPublish struct {
	ArticleId  uint     `json:"article_id" form:"article_id"` // 文章ID
	UserId     uint     `json:"user_id" form:"user_id"`
	Title      string   `json:"title" form:"title"`
	Status     string   `json:"status" form:"status"`           // 文章属性（草稿，私有，公开）
	CategoryID uint     `json:"category_id" form:"category_id"` // 所属类目ID，外键
	Summary    string   `json:"summary" form:"summary"`         // 文章摘要
	Content    string   `json:"content" form:"content"`         // 文章内容
	ImageUrl   string   `json:"image_url" form:"image_url"`     // 图片地址
	Tags       []string `json:"tags" form:"tags"`               // 标签
}

type TagsReq []struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}
type CategorysReq []struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}
