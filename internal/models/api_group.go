package models

import "gorm.io/gorm"

// ApiGroup
// @Description: api表和dict_item表的关联表, 用来记录api所属的分组
// @Author wangyulong 2024-10-19 10:22:19
type ApiGroup struct {
	gorm.Model
	ApiId   uint `json:"api_id"`
	GroupId uint `json:"group_id"`
}
