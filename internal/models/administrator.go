package models

import "gorm.io/gorm"

// Administrator 管理者
type Administrator struct {
	gorm.Model        //ID CreatedAt UpdatedAt DeletedAt
	Username   string `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Password   string `json:"password" gorm:"type:varchar(255);not null"`
}
