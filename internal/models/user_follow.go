package models

// UserFollow 用户关注
type UserFollow struct {
	FollowerId uint `json:"follower_id"` // 关注者
	FollowedId uint `json:"followed_id"` // 被关注者
}
