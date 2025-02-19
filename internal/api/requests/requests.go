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

// GetAllApiRes
// @Description: 获取所有api列表
// @Author wangyulong 2024-10-14 22:26:30
type GetAllApiRes struct {
	ApiList []*List `json:"api_list"`
}

type List struct {
	GroupID   uint        `json:"key"`
	GroupName string      `json:"title"`
	Children  []*Children `json:"children"`
}

type Children struct {
	ApiID   uint   `json:"key"`
	ApiName string `json:"title"`
}

type Result struct {
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
	Groups []*ApiGroup `json:"groups"`
}

type ApiGroup struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// ApiReqMethodRes 获取所有请求方法
type ApiReqMethodRes struct {
	Methods []*ApiReqMethod `json:"methods"`
}

type ApiReqMethod struct {
	Value string `json:"value"`
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
	BriefIntroduction string `json:"brief_introduction"` // API简介
	GroupId           uint   `json:"grouping_id"`        // API分组ID
	GroupName         string `json:"grouping"`           // API分组
	RequestMethodId   uint   `json:"request_method_id"`  //请求方式Id
	RequestMethodName string `json:"request_method"`     // 请求方式
}

// SearchApiListRes 检索api列表响应结构体
type SearchApiListRes struct {
	Api   []SearchApiRes `json:"api"`
	Total int            `json:"total"`
}

// MiddleResultG 用于获取所有分组的id和label
type MiddleResultG struct {
	ID    uint   `json:"id"`
	Label string `json:"label"`
}
