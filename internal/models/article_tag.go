package models

import "gorm.io/gorm"

type ArticleTag struct {
	gorm.Model
	ArticleId int `json:"article_id"` // 文章id
	TagId     int `json:"tag_id"`     // 标签id
}