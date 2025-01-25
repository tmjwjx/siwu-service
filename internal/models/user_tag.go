package models

import "gorm.io/gorm"

// UserTag 用户标签
type UserTag struct {
	gorm.Model
	UserID uint `json:"user_id"` // 用户id
	TagID  uint `json:"tag_id"`  // 标签id
}