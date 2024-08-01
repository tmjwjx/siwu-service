package models

// Administrator 管理者
type Administrator struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Password string `json:"password" gorm:"type:varchar(255);not null"`
}
