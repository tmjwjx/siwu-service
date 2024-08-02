package repository

import (
	"forum/internal/image/request"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

func InsertFile(attachment *request.Attachment) *gorm.DB {
	// 向数据库中存入文件数据
	result := globals.DB.Create(attachment)
	return result
}
