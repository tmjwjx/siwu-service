package requests

// RegisterReq 注册请求
type RegisterReq struct {
	Email      string `json:"email"`       // 邮箱
	VerifyCode string `json:"verify_code"` // 验证码
	Password   string `json:"password"`    // 密码
	RePassword string `json:"re_password"` // 重复密码
}

// LogicReq 登录请求
type LogicReq struct {
	Email    string `json:"email"`    // 邮箱
	Password string `json:"password"` // 密码
}

// // VerifyCodeReq 验证码消息请求
// type VerifyCodeReq struct {
// 	Email string `json:"email" form:"email"` // 邮箱
// }

// FollowReq 关注，取消关注请求
type FollowReq struct {
	FollowerId uint `json:"follower_id"`
	FollowedId uint `json:"followed_id"`
}

// UserRankReq 用户排行请求
type UserRankReq struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// UserRankRes 用户排行响应
type UserRankRes struct {
	Id              uint   `json:"id"`
	Nickname        string `json:"nickname"`
	CareerDirection string `json:"career_direction"`
	AvatarPath      string `json:"avatar_path"` // 头像图片
	IsFollowed      int    `json:"is_followed"` // 用户是否关注了这个排行榜上的用户：未关注：0，已关注：1，这个用户是自己：2
}
