package models

type UserFollow struct {
	FollowedId int `json:"followed_id"`
	FollowerId int `json:"follower_id"`
}