package models

// Attachment 图片附件表
type Attachment struct {
	Id          uint    `json:"id" gorm:"primaryKey"` // 主键
	FileName    string  `json:"file_name"`            // 文件名
	FileType    string  `json:"file_type"`            // 文件类型，例如 image/png
	FileSize    int64   `json:"file_size"`            // 文件大小（以字节为单位）
	FileContent []byte  `json:"file_content"`
	ArticleID   uint    `json:"article_id"` // 所属文章 ID，外键
	Article     Article `gorm:"foreignKey:ArticleID"`
}
