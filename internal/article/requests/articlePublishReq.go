package requests

type ReqPublish struct {
	UserId     uint     `json:"user_id" form:"user_id"`
	Title      string   `json:"title" form:"title"`
	Status     string   `json:"status" form:"status"`           // 文章属性（草稿，私有，公开）
	CategoryID uint     `json:"category_id" form:"category_id"` // 所属类目ID，外键
	Summary    string   `json:"summary" form:"summary"`         // 文章摘要
	Content    string   `json:"content" form:"content"`         // 文章内容
	Tag        []string `json:"tag" form:"tag"`                 // 标签
}