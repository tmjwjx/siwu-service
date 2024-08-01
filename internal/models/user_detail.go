package models

import (
	"gorm.io/gorm"
	"time"
)

// UserDetail 用户详细信息
type UserDetail struct {
	gorm.Model                   //ID CreatedAt UpdatedAt DeletedAt
	UserID            uint       `json:"user_id" gorm:"unique"` // 用户ID，外键，唯一索引
	CareerDirection   string     `json:"career_direction"`      // 职业方向
	UserHomePage      string     `json:"user_home_page"`        // 个人主页
	UserSignature     string     `json:"user_signature"`        // 个人签名
	UserTags          string     `json:"user_tags"`             // 作者标签
	BlogLink          string     `json:"blog_link"`             // 个人博客链接
	WeiboLink         string     `json:"weibo_link"`            // 新浪微博链接
	GithubLink        string     `json:"github_link"`           // Github链接
	LastLoginTime     *time.Time `json:"last_login_time"`       // 上次登录时间
	CreateAccountTime *time.Time `json:"create_account_time"`   // 创建账号时间
}
