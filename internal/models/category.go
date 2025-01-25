package models

import "gorm.io/gorm"

// Category
//
//	@Description: 类目表
type Category struct {
	gorm.Model        //ID CreatedAt UpdatedAt DeletedAt
	Name       string `json:"name" gorm:"not null;unique"` // 类目名称
	Icon       string `json:"icon_path" gorm:"not null"`   // 类目图标（直接存储图片路径）
	Articles   []Article
}