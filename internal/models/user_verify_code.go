package models

import "gorm.io/gorm"

// // UserVerifyCode 用户验证码信息
// type UserVerifyCode struct {
// 	gorm.Model        // ID CreatedAt UpdatedAt DeletedAt
// 	UserID     uint   `json:"user_id"`     // 用户ID
// 	VerifyCode string `json:"verify_code"` // 验证码，唯一
// }

// UserVerifyCode 用户验证码信息
type UserVerifyCode struct {
	gorm.Model        // ID CreatedAt UpdatedAt DeletedAt
	Email      string `json:"email"`       // 用户ID
	VerifyCode string `json:"verify_code"` // 验证码，唯一
}
