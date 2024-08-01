package models

// Follow 关注关系表
type Follow struct {
	ID         uint `json:"id" gorm:"primaryKey"`     // 主键
	FollowerID uint `json:"follower_id" gorm:"index"` // 关注者ID
	FollowedID uint `json:"followed_id" gorm:"index"` // 被关注者ID

	// FollowerID uint `json:"follower_id" gorm:"foreignKey:FollowerID;references:ID"`  // 关注者ID
	// FollowedID uint `json:"followed_id"  gorm:"foreignKey:FollowedID;references:ID"` // 被关注者ID
	User1 User `gorm:"foreignKey:FollowerID;references:ID"`
	User2 User `gorm:"foreignKey:FollowedID;references:ID"`
}
