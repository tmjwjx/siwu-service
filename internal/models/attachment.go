package models

import "gorm.io/gorm"

// Attachment 图片附件表
type Attachment struct {
	gorm.Model        //ID CreatedAt UpdatedAt DeletedAt
	Home       string `json:"home"`                 // 图片所属单位，即属于文章图片还是用户图片或是资源表图片
	HomeID     uint   `json:"homeID"`               // 图片对应的具体文章或用户的ID
	Type       string `json:"type"`                 // 文件类型，例如 images/png
	Size       int    `json:"size"`                 // 文件大小（以字节为单位）
	Path       string `json:"path" gorm:"not null"` // 文件路径
}
