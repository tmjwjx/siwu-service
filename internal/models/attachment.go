package models

import "gorm.io/gorm"

// Attachment 图片附件表
type Attachment struct {
	gorm.Model        //ID CreatedAt UpdatedAt DeletedAt
	Flag       int    `json:"flag"` // 图片所属单位的标志
	Type       string `json:"type"` // 文件类型，例如 image/png
	Size       int    `json:"size"` // 文件大小（以字节为单位）
	Path       string `json:"path" gorm:"not null"` // 文件路径
}