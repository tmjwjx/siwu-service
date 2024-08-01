package models

// UserFollowsTag 中间表
type UserFollowsTag struct {
	// ID     uint `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID uint `json:"user_id" gorm:"uniqueIndex:idx_user_tag;not null;foreignKey:UserID;references:ID"`
	TagID  uint `json:"tag_id" gorm:"uniqueIndex:idx_user_tag;not null;foreignKey:TagID;references:ID"`
	// CreatedAt time.Time `json:"created_at"` // 记录关注时间
}
