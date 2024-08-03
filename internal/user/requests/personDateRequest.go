package requests

import (
	"gorm.io/gorm"
	"time"
)

// User 用户简略信息
type User struct {
	gorm.Model        //ID CreatedAt UpdatedAt DeletedAt
	Nickname   string `form:"nickname" validate:"nickname"` // 用户名	1
	Email      string `form:"email" validate:"email"`       // 邮箱，唯一	2
	Password   string `form:"password" validate:"password"` // 密码	2
	Heat       int    `form:"heat"`                         // 个人热度
	FansCount  uint   `form:"fans_count"`                   // 粉丝数
	Path       string `form:"path"`                         // 用户头像
	UserDetail UserDetail
}

// UserDetail 用户详细信息
type UserDetail struct {
	gorm.Model                   //ID CreatedAt UpdatedAt DeletedAt
	UserID            uint       `form:"user_id"`             // 用户ID，外键，唯一索引
	CareerDirection   string     `form:"career_direction"`    // 职业方向	1
	UserHomePage      string     `form:"user_home_page"`      // 个人主页	1
	UserSignature     string     `form:"user_signature"`      // 个人签名	1
	UserTags          string     `form:"user_tags"`           // 作者标签	1
	BlogLink          string     `form:"blog_link"`           // 个人博客链接
	WeiboLink         string     `form:"weibo_link"`          // 新浪微博链接	2
	GithubLink        string     `form:"github_link"`         // Github链接	2
	LastLoginTime     *time.Time `form:"last_login_time"`     // 上次登录时间
	CreateAccountTime *time.Time `form:"create_account_time"` // 创建账号时间
}
