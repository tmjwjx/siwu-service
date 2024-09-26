package models

import "gorm.io/gorm"

// Resource 资源表
type Resource struct {
	gorm.Model        //ID CreatedAt UpdatedAt DeletedAt
	Icon       string `json:"icon_path" gorm:"not null"`   // 网站图标(直接存储图片路径)
	HelloWorld string `json:"hello_world" gorm:"not null"` //标题
}
