package models

// ArticleLike 文章点赞表
type ArticleLike struct {
	ID        uint    `json:"id" gorm:"primaryKey"`    // 主键
	ArticleID uint    `json:"article_id" gorm:"index"` // 文章ID
	UserID    uint    `json:"user_id" gorm:"index"`    // 点赞人ID
	Article   Article `gorm:"foreignKey:ArticleID"`
}
