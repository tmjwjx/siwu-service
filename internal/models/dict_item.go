package models

import "gorm.io/gorm"

// DictItem 字典项表
type DictItem struct {
	gorm.Model
	DictTypeCode string `json:"dict_type_code" gorm:"size:32"` // 外键关联字典类型（和字典类型编码关联）
	Label        string `json:"label" gorm:"size:50"`          // 字典项标签，字典项的显示名称（如“男”）
	Value        int    `json:"value" gorm:"not null"`         // 字典项键值，字典项的值（如 1 代表男）
	Sort         int    `json:"sort" gorm:"default:0"`         // 字典项排序，定义字典项的显示顺序
	Status       int    `json:"status" gorm:"default:1"`       // 字典项状态。（如：启用、禁用）
	Description  string `json:"description" gorm:"size:100"`   // 字典项描述
	ExtendValue  string `json:"extend_value" gorm:"size:100"`  // 扩展值
}
