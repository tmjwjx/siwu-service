package models

import "gorm.io/gorm"

type SystemMessage struct {
	gorm.Model
	ReceiverID uint                   `json:"receiver_id"`              // 接收者id
	SenderID   uint                   `json:"sender_id"`                // 发送者id
	Type       string                 `json:"type"`                     // 消息类型
	Content    map[string]interface{} `json:"content" gorm:"type:json"` // JSON格式 消息内容
}