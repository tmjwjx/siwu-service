package requests

type GetArticleInfoFirstReq struct {
	Content   string   `json:"content"`    // 文章的全部内容
	Tags      []string `json:"tags"`       // 所有标签, 用于给文章匹配相应的标签
	ArticleID uint     `json:"article_id"` // 文章ID
}

type GetArticleInfoFirstRes struct {
	Key      string   `json:"key"`      // hash值
	Abstract string   `json:"abstract"` // 文章的摘要
	Summary  string   `json:"summary"`  // 文章的总结
	Tags     []string `json:"tags"`     // 与文章相匹配的标签
}

type SaveArticleIDReq struct {
	Key       string `json:"key"`        // hash值
	ArticleID uint   `json:"article_id"` // 文章ID
}

type GetArticleInfoReq struct {
	ArticleID uint `json:"article_id"` // 文章ID
	UserID    uint `json:"user_id"`    // 用户ID
}

type GetArticleInfoRes struct {
	Abstract string             `json:"abstract"` // 文章的摘要
	Summary  string             `json:"summary"`  // 文章的总结
	Codes    []*CodeExplanation `json:"codes"`
}

type CodeExplanation struct {
	Question    string `json:"question"`    // 代码提问
	Explanation string `json:"explanation"` // 代码解释
}

type DelArticleInfoReq struct {
	ArticleID uint `json:"article_id"` // 文章ID
}
