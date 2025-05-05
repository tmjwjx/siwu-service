package models

import (
	"forum/pkg/globals"
	"gorm.io/gorm"
)

// Attachment 图片附件表
type Attachment struct {
	gorm.Model              //ID CreatedAt UpdatedAt DeletedAt
	Home       globals.Home `json:"home"`    // 图片所属单位，即属于文章图片还是用户图片或是资源表图片
	HomeID     uint         `json:"home_id"` // 图片对应的具体文章或用户的ID
	Name       string       `json:"name"`    // 文件名
	Type       string       `json:"type"`    // 文件类型，例如 images/png
	Size       int64        `json:"size"`    // 文件大小（以字节为单位）
	Path       string       `json:"path"`    // 文件路径
}
