package requests

type ReqSearch struct {
	UserId     int    `json:"user_id" form:"user_id"`         // 用户id
	Keyword    string `json:"keyword" form:"keyword"`         // 搜索词
	CategoryId int    `json:"category_id" form:"category_id"` // 类目
	Page       int    `json:"page" form:"page"`               // 页数
	Limit      int    `json:"limit" form:"limit"`             // 每页条数
	Kind       int    `json:"kind" form:"kind"`               // 排序
}