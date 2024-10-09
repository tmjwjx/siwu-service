package requests

type SystemRes struct {
	UserId    uint   `json:"user_id"`
	ArticleId uint   `json:"article_id"`
	CreatedAt string `json:"created_at"`
}