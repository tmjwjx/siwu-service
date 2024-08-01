package models

// Resource 资源表
type Resource struct {
	ID            uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Logo          []byte `json:"logo" gorm:"type:longblob"`
	HelloWorld    string `json:"hello_world" gorm:"type:varchar(255);not null"`
	Advertisement []byte `json:"advertisement" gorm:"type:longblob"`
}
