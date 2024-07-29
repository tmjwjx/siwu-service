package models

import "time"

// Tag 标签
type Tag struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string    `json:"name" gorm:"type:varchar(100);not null;unique"`
	Description   string    `json:"description" gorm:"type:text"`
	ArticleCount  int       `json:"articleCount" gorm:"default:0"`
	Popularity    int       `json:"popularity" gorm:"default:0"`
	FollowerCount int       `json:"followerCount" gorm:"default:0"`
	Image         []byte    `json:"image" gorm:"type:longblob"`
	Articles      []Article `json:"articles" gorm:"many2many:article_tags;"`
	Followers     []User    `json:"followers" gorm:"many2many:user_follows_tags;"`
}

// Resource 资源表
type Resource struct {
	ID            uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Logo          []byte `json:"logo" gorm:"type:longblob"`
	HelloWorld    string `json:"helloWorld" gorm:"type:varchar(255);not null"`
	Advertisement []byte `json:"advertisement" gorm:"type:longblob"`
}

// Administrator 管理者
type Administrator struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Password string `json:"password" gorm:"type:varchar(255);not null"`
}

// ArticleTag 中间表
type ArticleTag struct {
	ID        uint `json:"id" gorm:"primaryKey;autoIncrement"`
	ArticleID uint `json:"articleID" gorm:"index"`
	TagID     uint `json:"tagID" gorm:"index"`
}

// UserFollowsTag 中间表
type UserFollowsTag struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint      `json:"userID" gorm:"index"`
	TagID     uint      `json:"tagID" gorm:"index"`
	CreatedAt time.Time `json:"createdAt"` // 记录关注时间
}
