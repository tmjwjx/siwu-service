package models

// ArticleTag 中间表
type ArticleTag struct {
	// ID        uint `json:"id" gorm:"primaryKey;autoIncrement"`
	ArticleID uint `json:"article_id" gorm:"uniqueIndex:idx_article_tag;foreignKey:ArticleID;references:ID"`
	TagID     uint `json:"tag_id" gorm:"uniqueIndex:idx_article_tag;foreignKey:TagID;references:ID"`
}
