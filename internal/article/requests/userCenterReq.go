package requests

// UserArticleOrCollectionReq
// @Description: 用户文章或收藏列表请求
// @Author tianjiajie 2024-10-18 15:36:06
type UserArticleOrCollectionReq struct {
	Id      int    `json:"id" form:"id"`
	Type    string `json:"type" form:"type"`
	Page    int    `json:"page" form:"page"`
	Limit   int    `json:"limit" form:"limit"`
	Keyword string `json:"keyword" form:"keyword"`
}