package requests

import "forum/internal/models"

/*// User 用户简略信息
type User struct {
	gorm.Model        //ID CreatedAt UpdatedAt DeletedAt
	Nickname   string `json:"nickname" form:"nickname" validate:"nickname"` // 用户名	1
	Email      string `json:"email" form:"email" validate:"email"`          // 邮箱，唯一	2
	Password   string `json:"password" form:"password" validate:"password"` // 密码	2
	Heat       int    `json:"heat" form:"heat"`                             // 个人热度
	FansCount  uint   `json:"fans_count" form:"fans_count"`                 // 粉丝数
	Path       string `json:"path" form:"path"`                             // 用户头像
	UserDetail UserDetail
}

// UserDetail 用户详细信息
type UserDetail struct {
	gorm.Model                   //ID CreatedAt UpdatedAt DeletedAt
	UserID            uint       `json:"user_id" form:"user_id"`                         // 用户ID，外键，唯一索引
	CareerDirection   string     `json:"career_direction" form:"career_direction"`       // 职业方向	1
	UserHomePage      string     `json:"user_home-page" form:"user_home_page"`           // 个人主页	1
	UserSignature     string     `json:"user_signature" form:"user_signature"`           // 个人签名	1
	UserTags          string     `json:"user_tags" form:"user_tags"`                     // 作者标签	1
	BlogLink          string     `json:"blog_link" form:"blog_link"`                     // 个人博客链接
	WeiboLink         string     `json:"weibo_link" form:"weibo_link"`                   // 新浪微博链接	2
	GithubLink        string     `json:"github_link" form:"github_link"`                 // Github链接	2
	LastLoginTime     *time.Time `json:"last_login_time" form:"last_login_time"`         // 上次登录时间
	CreateAccountTime *time.Time `json:"create_account_time" form:"create_account_time"` // 创建账号时间
}*/

type UserResponse struct {
	User        models.User       `json:"user"`
	UserDetail  models.UserDetail `json:"user_detail"`
	UserTags    []string          `json:"user_tags"`
	AllTagNames []string          `json:"all_tag_names"`
	Path        string            `json:"path"`
}
