package requests

// BackstageLoginReq 后台登录请求
type BackstageLoginReq struct {
	Email    string `json:"email"`    // 邮箱
	Password string `json:"password"` // 密码
}

// BackstageLoginRes 后台登录响应
type BackstageLoginRes struct {
	Id         uint     `json:"id"`          // 用户id
	Nickname   string   `json:"nickname"`    // 用户名
	Email      string   `json:"email"`       // 邮箱
	UserStatus int      `json:"user_status"` // 用户状态
	RoleNames  []string `json:"role_names"`  // 角色的名称集合
	AvatarPath string   `json:"avatar_path"` // 用户头像
	Code       string   `json:"code"`        // 当前用户所拥有的排序最高的角色编码
}
