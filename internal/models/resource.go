package models

import "gorm.io/gorm"

// Resource 资源表
type Resource struct {
	gorm.Model           //ID CreatedAt UpdatedAt DeletedAt
	Logo          []byte `json:"logo" gorm:"type:longblob"`                     // 网站图标
	HelloWorld    string `json:"hello_world" gorm:"type:varchar(255);not null"` //标题
	Advertisement []byte `json:"advertisement" gorm:"type:longblob"`            // 广告
}
