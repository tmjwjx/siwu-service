package requests

// UserDataReq 用户个人资料更新
type UserDataReq struct {
	ID              uint     `json:"id" form:"id"`                             // 用户ID
	Nickname        string   `json:"nickname" form:"nickname"`                 // 用户名
	CareerDirection string   `json:"career_direction" form:"career_direction"` // 职业方向
	HomePage        string   `json:"user_home-page" form:"user_home_page"`     // 个人主页
	Signature       string   `json:"user_signature" form:"user_signature"`     // 个人签名
	Path            string   `json:"path" form:"path"`                         // 用户头像
	UserTags        []string `json:"user_tags" form:"user_tags"`               // 用户标签
}

// UserDataRes 用户个人资料响应
type UserDataRes struct {
	ID              uint     `json:"id"`               // 用户ID
	Nickname        string   `json:"nickname"`         // 用户名
	CareerDirection string   `json:"career_direction"` // 职业方向
	HomePage        string   `json:"user_home_page"`   // 个人主页
	Signature       string   `json:"user_signature"`   // 个人签名
	Path            string   `json:"path"`             // 用户头像
	UserTags        []string `json:"user_tags"`        // 用户标签
	AllTagNames     []string `json:"all_tag_names"`    // 所有标签的名字
}

// UserAccountReq 用户账号设置更新
type UserAccountReq struct {
	ID         uint   `json:"id"`          // 用户ID
	Email      string `json:"email"`       // 邮箱，唯一
	BlogLink   string `json:"blog_link"`   // 个人博客链接
	WeiboLink  string `json:"weibo_link"`  // 新浪微博链接
	GithubLink string `json:"github_link"` // Github链接
	Password   string `json:"password"`    // 密码
}

// UserAccountRes 用户账号设置响应
type UserAccountRes struct {
	ID         uint   `json:"id"`          // 用户ID
	Email      string `json:"email"`       // 邮箱，唯一
	BlogLink   string `json:"blog_link"`   // 个人博客链接
	WeiboLink  string `json:"weibo_link"`  // 新浪微博链接
	GithubLink string `json:"github_link"` // Github链接
	Password   string `json:"password"`    // 密码
}

// UserPrivateSettingsReq 用户私信设置更新
type UserPrivateSettingsReq struct {
	ID              uint   `json:"id"`
	PrivateSettings string `json:"private_settings"`
}

// UserPrivateSettingsRes 用户私信设置响应
type UserPrivateSettingsRes struct {
	ID              uint   `json:"id"`
	PrivateSettings string `json:"private_settings"`
}
