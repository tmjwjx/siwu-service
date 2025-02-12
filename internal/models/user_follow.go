package models

import "gorm.io/gorm"

// UserFollow 用户关注
type UserFollow struct {
	gorm.Model      // ID CreatedAt UpdatedAt DeletedAt
	FollowerId uint `json:"follower_id"`                  // 关注者
	FollowedId uint `json:"followed_id"`                  // 被关注者
	IsRead     bool `json:"is_read" gorm:"default:false"` // 是否已读
}
