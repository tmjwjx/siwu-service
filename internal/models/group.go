package models

import "gorm.io/gorm"

// Group API分组
type Group struct {
	gorm.Model
	Name string `json:"name"` // 分组名称
}
