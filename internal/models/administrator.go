package models

import (
	"gorm.io/gorm"
	"time"
)

// Administrator 管理员
type Administrator struct {
	gorm.Model              // ID CreatedAt UpdatedAt DeletedAt
	LastLoginTime time.Time `json:"last_login_time" gorm:"default:'0001-01-01 00:00:00'"` // 最后一次的登录时间
	Email         string    `json:"email" gorm:"not null;unique"`                         // 邮箱，唯一
	Name          string    `json:"name" gorm:"size:100"`                                 // 名称
	Password      string    `json:"password" gorm:"size:100;not null"`                    // 密码
	// 头像
}
