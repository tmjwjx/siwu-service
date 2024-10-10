package requests

// TagFansCountRes 更新人数响应
type TagFansCountRes struct {
	TagFansCount *TagFansCount `json:"tag_fans_count"`
}

type TagFansCount struct {
	FansCount string `json:"fans_count"`
}

// TagReq 请求标签结构体
type TagReq struct {
	ID uint `json:"id" form:"id"`
}

// TagRes 响应标签结构体
type TagRes struct {
	TagList []*Tag `json:"tag_list"`
}

type Tag struct {
	ID           uint   `json:"id" form:"id"`                       // 标签ID
	Name         string `json:"name" form:"name"`                   // 标签名称
	Description  string `json:"description" form:"description"`     // 标签描述
	ArticleCount string `json:"article_count" form:"article_count"` // 标签关联的文章数量
	Heat         string `json:"heat" form:"heat"`                   // 标签热度
	FansCount    string `json:"fans_count" form:"fans_count"`       // 关注人数
	Path         string `json:"path" form:"path"`                   // 标签头像
}

// BsAddTagReq 后台新增请求
type BsAddTagReq struct {
	Name         string   `json:"name" form:"name"`                   // 标签名称
	Description  string   `json:"description" form:"description"`     // 标签描述
	ArticleCount int      `json:"article_count" form:"article_count"` // 标签关联的文章数量
	Heat         int      `json:"heat" form:"heat"`                   // 标签热度
	FansCount    int      `json:"fans_count" form:"fans_count"`       // 关注人数
	Path         []string `json:"path" form:"path"`                   // 标签头像
}

// BsDelTagReq 后台删除请求
type BsDelTagReq struct {
	ID uint `json:"id"` // 标签ID
}

// BsBatchDelTagReq 后台批量删除请求
type BsBatchDelTagReq struct {
	ID []uint `json:"id"` // 标签ID
}

// BsUpTagReq 后台更新请求
type BsUpTagReq struct {
	ID           uint     `json:"id" form:"id"`                       // 标签ID
	Name         string   `json:"name" form:"name"`                   // 标签名称
	Description  string   `json:"description" form:"description"`     // 标签描述
	ArticleCount int      `json:"article_count" form:"article_count"` // 标签关联的文章数量
	Heat         int      `json:"heat" form:"heat"`                   // 标签热度
	FansCount    int      `json:"fans_count" form:"fans_count"`       // 关注人数
	Path         []string `json:"path" form:"path"`                   // 标签头像
}

// BsQueTagReq 后台查询请求
type BsQueTagReq struct {
	Name string `json:"name"` // 标签名称
}

// BsQueTagRes 后台查询响应
type BsQueTagRes struct {
	ID           uint   `json:"id"`            // 标签ID
	Name         string `json:"name"`          // 标签名称
	Description  string `json:"description"`   // 标签描述
	ArticleCount int    `json:"article_count"` // 标签关联的文章数量
	Heat         int    `json:"heat"`          // 标签热度
	FansCount    int    `json:"fans_count"`    // 关注人数
	Path         string `json:"path"`          // 标签头像
}

// BsBatchQueTagReq 后台批量查询请求
type BsBatchQueTagReq struct {
	Offset int `json:"offset"` // 分页查询的起始位置
	Limit  int `json:"limit"`  // 分页查询返回的记录条数
}
