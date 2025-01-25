package models

// ApiGroup
// @Description: api表和group表的关联表
// @Author wangyulong 2024-10-19 10:22:19
type ApiGroup struct {
	ApiId   uint `json:"api_id"`
	GroupId uint `json:"group_id"`
}
