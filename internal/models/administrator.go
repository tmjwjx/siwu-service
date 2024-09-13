package models

import "gorm.io/gorm"

// Administrator 管理者
type Administrator struct {
	gorm.Model        //ID CreatedAt UpdatedAt DeletedAt
	Email      string `json:"username" gorm:"size:16;not null"`
	Password   string `json:"password" gorm:"size:16;not null"`
}