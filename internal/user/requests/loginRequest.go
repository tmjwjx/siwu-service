package requests

// RegisterMsg 注册消息
type RegisterMsg struct {
	Email      string `json:"email"`
	Code       string `json:"code"`
	Password   string `json:"password"`
	RePassword string `json:"re_password"`
}
