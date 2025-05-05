package models

import "gorm.io/gorm"

// Api 管理表
type Api struct {
	gorm.Model
	Path              string `json:"path"`               // API路径
	BriefIntroduction string `json:"brief_introduction"` // API简介
}
