package models

import "gorm.io/gorm"

// UserMessage 用户消息表(小铃铛专属)
type UserMessage struct {
	gorm.Model            //ID CreatedAt UpdatedAt DeletedAt
	UserID           uint `json:"user_id" gorm:"unique;index"`        // 用户ID，外键，唯一索引
	UnreadComment    int  `json:"unread_comment" gorm:"default:0"`    // 未读评论数量
	UnreadLike       int  `json:"unread_like" gorm:"default:0"`       // 未读点赞数量
	UnreadCollection int  `json:"unread_collection" gorm:"default:0"` // 未读收藏数量
	UnreadFan        int  `json:"unread_fan" gorm:"default:0"`        // 未读粉丝数量
	UnreadMessage    int  `json:"unread_message" gorm:"default:0"`    // 未读私信数量
	UnreadSystem     int  `json:"unread_system" gorm:"default:0"`     // 未读系统消息数量
}