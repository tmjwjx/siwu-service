package models

import "gorm.io/gorm"

// Administrator 管理员。
type Administrator struct {
	gorm.Model        // ID CreatedAt UpdatedAt DeletedAt
	Email      string `json:"email" gorm:"not null;unique"`      // 邮箱，唯一
	Password   string `json:"password" gorm:"size:100;not null"` // 密码
}
