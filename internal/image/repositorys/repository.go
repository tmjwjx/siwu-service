package repositorys

import (
	"forum/internal/image/requests"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

func InsertFile(attachment *requests.Attachment) *gorm.DB {
	// 向数据库中存入文件数据
	result := globals.DB.Create(attachment)
	return result
}
