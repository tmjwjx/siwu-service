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

// Resource 资源表
type Resource struct {
	ID            uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Logo          []byte `json:"logo" gorm:"type:longblob"`
	HelloWorld    string `json:"hello_world" gorm:"type:varchar(255);not null"`
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
	//ID        uint `json:"id" gorm:"primaryKey;autoIncrement"`
	ArticleID uint `json:"article_id" gorm:"uniqueIndex:idx_article_tag;foreignKey:ArticleID;references:ID"`
	TagID     uint `json:"tag_id" gorm:"uniqueIndex:idx_article_tag;foreignKey:TagID;references:ID"`
}

// UserFollowsTag 中间表
type UserFollowsTag struct {
	//ID     uint `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID uint `json:"user_id" gorm:"uniqueIndex:idx_user_tag;not null;foreignKey:UserID;references:ID"`
	TagID  uint `json:"tag_id" gorm:"uniqueIndex:idx_user_tag;not null;foreignKey:TagID;references:ID"`
	//CreatedAt time.Time `json:"created_at"` // 记录关注时间
}
