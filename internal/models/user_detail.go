package models

import (
	"gorm.io/gorm"
)

// UserDetail 用户详细信息
type UserDetail struct {
	gorm.Model             // ID CreatedAt UpdatedAt DeletedAt
	UserID          uint   `json:"user_id" gorm:"unique;index"` // 用户ID，外键，唯一索引
	CareerDirection string `json:"career_direction"`            // 职业方向
	HomePage        string `json:"home_page"`                   // 个人主页
	Signature       string `json:"signature"`                   // 个人签名
	BlogLink        string `json:"blog_link"`                   // 个人博客链接
	WeiboLink       string `json:"weibo_link"`                  // 新浪微博链接
	GithubLink      string `json:"github_link"`                 // Github链接
}
