package models

import (
	"gorm.io/gorm"
)

// Role 角色
type Role struct {
	gorm.Model
	Name   string `json:"name" gorm:"size:16;not null"` // 名称
	Code   string `json:"code" gorm:"size:16;not null"` // 编码（名称的英文单词）
	Status int    `json:"status" gorm:"default:0"`      // 状态 0全部 1正常 2封禁（该角色是否可用）
	Sort   int    `json:"sort" gorm:"default:0"`        // 用于角色排序（角色在前端的显示规则）
}
