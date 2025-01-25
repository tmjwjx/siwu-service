package models

import "gorm.io/gorm"

// Tag 标签
type Tag struct {
	gorm.Model             //ID CreatedAt UpdatedAt DeletedAt
	Name         string    `json:"name" gorm:"size:16;not null"`   // 标签名称
	Description  string    `json:"description" gorm:"type:text"`   // 标签描述
	ArticleCount int       `json:"article_count" gorm:"default:0"` // 标签关联的文章数量
	Heat         int       `json:"heat" gorm:"default:0"`          // 标签热度
	FansCount    int       `json:"fans_count" gorm:"default:0"`    // 关注人数
	Articles     []Article `gorm:"many2many:article_tags;"`        // 文章和标签 多对多
	Users        []User    `gorm:"many2many:user_tags"`            // 用户和标签 多对多
}