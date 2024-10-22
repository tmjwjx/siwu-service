package models

import "gorm.io/gorm"

// RequestMethod
// @Description: api的请求方法表
// @Author wangyulong 2024-10-19 10:34:26
type RequestMethod struct {
	gorm.Model
	Name string `json:"name"` //请求方法名
}
