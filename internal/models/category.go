package models

import "gorm.io/gorm"

// Category 类目表
type Category struct {
	gorm.Model        //ID CreatedAt UpdatedAt DeletedAt
	Name       string `json:"name" gorm:"not null;unique"` // 类目名称
	//Path       string `json:"path" gorm:"not null"`        // 类目图标
	Articles []Article
}
