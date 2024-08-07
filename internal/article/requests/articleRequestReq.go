package requests

type ReqSearch struct {
	UserId     int    `json:"user_id" form:"user_id"`         // 用户id
	Query      string `json:"query" form:"query"`             // 搜索词
	CategoryId int    `json:"category_id" form:"category_id"` // 类目
	Tag        string `json:"tag" form:"tag"`                 // 标签
	Page       int    `json:"page" form:"page"`               // 页数
	Limit      int    `json:"limit" form:"limit"`             // 每页条数
	Kind       int    `json:"kind" form:"kind"`               // 排序
}