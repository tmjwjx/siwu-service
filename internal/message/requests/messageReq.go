package requests

type MessageReq struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

type LikeMessageRes struct {
	UserId    uint   `json:"user_id"`
	Nickname  string `json:"nickname"`
	Path      string `json:"path"`
	ArticleId uint   `json:"article_id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}