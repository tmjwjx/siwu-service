package models

import "gorm.io/gorm"

// UserVerifyCode 用户验证码信息
type UserVerifyCode struct {
	gorm.Model        // ID CreatedAt UpdatedAt DeletedAt
	UserID     uint   `json:"user_id" gorm:"unique;index"` // 用户ID，外键，唯一索引
	VerifyCode string `json:"verify_code"`                 // 验证码
}
