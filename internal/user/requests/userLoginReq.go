package requests

// RegisterMsg 注册消息
type RegisterMsg struct {
	Email      string `json:"email"`       // 邮箱
	Code       string `json:"code"`        // 验证码
	Password   string `json:"password"`    // 密码
	RePassword string `json:"re_password"` // 重复密码
}

// LogicMsg 登录消息
type LogicMsg struct {
	Email    string `json:"email"`    // 邮箱
	Password string `json:"password"` // 密码
}

// FindPassword 找回密码
type FindPassword struct {
	Email string `json:"email"` // 邮箱
	Code  string `json:"code"`  // 验证码

}
