package models

import "gorm.io/gorm"

// DictType 字典类型表
type DictType struct {
	gorm.Model
	Name        string `json:"name" gorm:"size:16;not null"`     // 字典类型名称，描述该字典的用途（例如：用户状态、商品分类等）
	Code        string `json:"code" gorm:"unique;index;size:32"` // 字典类型编码，唯一标识字典类型的代码
	Status      int    `json:"status" gorm:"default:1"`          // 字典类型状态（如：启用1、禁用2）。
	Description string `json:"description" gorm:"size:100"`      // 字典类型描述
}
