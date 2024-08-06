package models

import "gorm.io/gorm"

// Resource 资源表
type Resource struct {
	gorm.Model               //ID CreatedAt UpdatedAt DeletedAt
	IconPath          string `json:"icon_path" gorm:"not null"` // 网站图标
	HelloWorld        string `json:"hello_world" gorm:"not null"` //标题
	AdvertisementPath string `json:"advertisement_path" gorm:"not null"` // 广告
}