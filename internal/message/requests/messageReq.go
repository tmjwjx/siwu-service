package requests

// MessageReq
// @Description: 消息请求
// @Author tianjiajie 2024-10-05 17:28:11
type MessageReq struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

// LikeAndCollectionMessageRes
// @Description: 点赞消息响应
// @Author tianjiajie 2024-10-05 17:32:06
type LikeAndCollectionMessageRes struct {
	UserId    uint   `json:"user_id"`
	Nickname  string `json:"nickname"`
	Path      string `json:"path"`
	ArticleId uint   `json:"article_id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}