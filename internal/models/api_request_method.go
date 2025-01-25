package models

// ApiRequestMethod
// @Description: api和requestMethod的关联表
// @Author wangyulong 2024-10-19 09:55:35
type ApiRequestMethod struct {
	ApiId           uint `json:"api_id"`
	RequestMethodId uint `json:"request_method_id"`
}
