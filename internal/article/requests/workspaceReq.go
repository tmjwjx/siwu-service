package requests

// HotArticleRes
// @Description: 查询热门文章数据响应
// @Author tianjiajie 2025-01-16 14:38:25
type HotArticleRes struct {
	ID         uint   `json:"id"`          // 文章ID
	Title      string `json:"title"`       // 文章标题
	LikesCount string `json:"likes_count"` // 点赞数
	Increase   string `json:"increases"`   // 点赞涨幅
}

type WorkspaceData struct {
	ArticleTotal  string `json:"article_total"`  // 文章总数
	NewAdd        string `json:"new_add"`        // 新增文章数
	TodayViews    int64  `json:"today_views"`    // 今日浏览量
	TodayComments int64  `json:"today_comments"` // 今日评论数
}
