package requests

// RegisterMsg 注册消息
type RegisterMsg struct {
	Email      string `json:"email"`       // 邮箱
	VerifyCode string `json:"verify_code"` // 验证码
	Password   string `json:"password"`    // 密码
	RePassword string `json:"re_password"` // 重复密码
}

// LogicMsg 登录消息
type LogicMsg struct {
	Email    string `json:"email"`    // 邮箱
	Password string `json:"password"` // 密码
}

// VerifyCodeMsg 验证码消息
type VerifyCodeMsg struct {
	Email string `json:"email" form:"email"` // 邮箱
}

// FollowMsg 关注，取消关注消息
type FollowMsg struct {
	FollowerId uint `json:"follower_id"`
	FollowedId uint `json:"followed_id"`
}

// UserRankMsg 用户排行消息
type UserRankMsg struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// UserRankReq 响应用户排行消息
type UserRankReq struct {
	Id              uint   `json:"id"`
	Nickname        string `json:"nickname"`
	CareerDirection string `json:"career_direction"`
	AvatarPath      string `json:"avatar_path"` // 头像图片
	IsFollowed      int    `json:"is_followed"` // 用户是否关注了这个排行榜上的用户：未关注：0，已关注：1，这个用户是自己：2
}
