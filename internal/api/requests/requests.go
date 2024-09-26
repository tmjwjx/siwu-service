package requests

// ApiInitRes 获取所有api列表
type ApiInitRes struct {
	ID                uint   `json:"id"`                 // api的ID
	Path              string `json:"path"`               // API路径
	GroupingId        uint   `json:"grouping_id"`        // API分组ID
	Grouping          string `json:"grouping"`           // API分组
	BriefIntroduction string `json:"brief_introduction"` // API简介
	RequestMethodId   uint   `json:"request_method_id"`  //请求方式Id
	RequestMethod     string `json:"request_method"`     // 请求方式
}

// ApiDetailsRes 获取当前api详情
type ApiDetailsRes struct {
	ID                uint   `json:"id"`                 // api的ID
	Path              string `json:"path"`               // API路径
	GroupingId        uint   `json:"grouping_id"`        // API分组ID
	Grouping          string `json:"grouping"`           // API分组
	BriefIntroduction string `json:"brief_introduction"` // API简介
	RequestMethodId   uint   `json:"request_method_id"`  //请求方式 Id
	RequestMethod     string `json:"request_method"`     // 请求方式
}

// ApiGroupRes 获取所有api分组列表
type ApiGroupRes struct {
	Groups []*ApiGroup
}

type ApiGroup struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}

// ApiReqMethodRes 获取所有请求方法
type ApiReqMethodRes struct {
	Methods []*ApiReqMethod
}

type ApiReqMethod struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}

// DeleteApiReq 删除api
type DeleteApiReq struct {
	ID []uint `json:"ids"`
}

// UpdateApiReq 编辑api
type UpdateApiReq struct {
	ID                uint   `json:"id"`                 // api的ID
	Path              string `json:"path"`               // API路径
	RequestMethod     string `json:"request_method"`     //请求方式Id
	Grouping          string `json:"grouping"`           // API分组
	BriefIntroduction string `json:"brief_introduction"` // API简介
}

// CreateApiReq 添加api
type CreateApiReq struct {
	Path              string `json:"path"`               // API路径
	RequestMethod     string `json:"request_method"`     //请求方式Id
	Grouping          string `json:"grouping"`           // API分组
	BriefIntroduction string `json:"brief_introduction"` // API简介
}

// SearchApiListReq 检索api列表
type SearchApiListReq struct {
	Path              string `json:"path"`               // API路径
	RequestMethod     string `json:"request_method"`     //请求方式Id
	Grouping          string `json:"grouping"`           // API分组
	BriefIntroduction string `json:"brief_introduction"` // API简介
	Page              int    `json:"page"`               // 分页查询的起始位置
	Limit             int    `json:"limit"`              // 分页查询要返回记录的数量
}

// SearchApiRes 检索api列表
type SearchApiRes struct {
	ID                uint   `json:"id"`                 // api的ID
	Path              string `json:"path"`               // API路径
	GroupingId        uint   `json:"grouping_id"`        // API分组ID
	Grouping          string `json:"grouping"`           // API分组
	BriefIntroduction string `json:"brief_introduction"` // API简介
	RequestMethodId   uint   `json:"request_method_id"`  //请求方式Id
	RequestMethod     string `json:"request_method"`     // 请求方式
}

// SearchApiListRes 检索api列表响应结构体
type SearchApiListRes struct {
	Api []*SearchApiRes
}
