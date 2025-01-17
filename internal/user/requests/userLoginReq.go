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

// LogicRes 登录响应
type LogicRes struct {
	Id       uint   `json:"id"`
	Nickname string `json:"nickname"`
}

// ForgotPasswordReq 忘记验证码请求
type ForgotPasswordReq struct {
	Email      string `json:"email"`       // 邮箱
	VerifyCode string `json:"verify_code"` // 验证码
	Password   string `json:"password"`    // 密码
	RePassword string `json:"re_password"` // 重复密码
}
