package models

// Tag 标签
type Tag struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string    `json:"name" gorm:"type:varchar(100);not null;unique"`
	Description   string    `json:"description" gorm:"type:text"`
	ArticleCount  int       `json:"article_count" gorm:"default:0"`
	Popularity    int       `json:"popularity" gorm:"default:0"`
	FollowerCount int       `json:"follower_count" gorm:"default:0"`
	Image         []byte    `json:"image" gorm:"type:longblob"`
	Articles      []Article `json:"articles" gorm:"many2many:article_tags;foreignKey:ID;joinForeignKey:TagID;references:ID;joinReferences:ArticleID"`
	Followers     []User    `json:"followers" gorm:"many2many:user_follows_tags;foreignKey:ID;joinForeignKey:TagID;references:ID;joinReferences:UserID"`
}
