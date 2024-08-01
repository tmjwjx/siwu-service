package models

// ArticleCollection 文章收藏表
type ArticleCollection struct {
	ID        uint    `json:"id" gorm:"primaryKey"`    // 主键
	ArticleID uint    `json:"article_id" gorm:"index"` // 文章ID
	UserID    uint    `json:"user_id" gorm:"index"`    // 收藏人ID
	Article   Article `gorm:"foreignKey:ArticleID"`
	User      User    `gorm:"foreignKey:UserID"`
}
