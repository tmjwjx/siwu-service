package models

type UserFollow struct {
	FollowerId int `json:"follower_id"` // 关注者
	FollowedId int `json:"followed_id"` // 被关注者
}