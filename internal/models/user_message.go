package models

import "gorm.io/gorm"

// UserMessage 用户消息表(小铃铛专属)
type UserMessage struct {
	gorm.Model             //ID CreatedAt UpdatedAt DeletedAt
	UserID            uint `json:"user_id" gorm:"uniqueIndex"` // 用户ID，外键，唯一索引
	UnreadComments    int  `json:"unread_comments"`            // 未读评论数量
	UnreadLikes       int  `json:"unread_likes"`               // 未读点赞数量
	UnreadCollections int  `json:"unread_collections"`         // 未读收藏数量
	UnreadFans        int  `json:"unread_fans"`                // 未读粉丝数量
	UnreadMessages    int  `json:"unread_messages"`            // 未读私信数量
	UnreadSystem      int  `json:"unread_system"`              // 未读系统消息数量
}
