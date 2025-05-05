package models

import "gorm.io/gorm"

// ApiDictItemRequestMethod
// @Description: api和dict_item的关联表,用来记录api所属的请求方法
// @Author wangyulong 2024-10-19 09:55:35
type ApiDictItemRequestMethod struct {
	gorm.Model
	ApiId           uint `json:"api_id"`
	RequestMethodId uint `json:"request_method_id"`
}
