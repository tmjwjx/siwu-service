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

// ReqVerifyCode 请求验证码
type ReqVerifyCode struct {
	Email string `json:"email" form:"email"` // 邮箱
}
