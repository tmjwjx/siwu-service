package models

import (
	"forum/pkg/globals"
	"gorm.io/gorm"
)

// Advertisement
// @Description:广告
type Advertisement struct {
	gorm.Model              //ID CreatedAt UpdatedAt DeletedAt
	Home       globals.Home `json:"home"`   // 广告图片属于哪个页面
	Status     int          `json:"status"` // 该广告是否被展示(1:展示 2:不展示)
	Path       string       `json:"path"`   // 广告图片路径
}
