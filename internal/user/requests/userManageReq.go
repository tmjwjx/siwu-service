package requests

type ReseatReq struct {
	Id uint `json:"id"`
}

// AddReq 添加用户请求
type AddReq struct {
	NickName   string `json:"nickname"`
	Email      string `json:"email"`
	UserStatus int    `json:"user_status"` // 用户状态 0全部 1正常 2封禁
	RoleIds    []uint `json:"role_ids"`
	AvatarPath string `json:"avatar_path"`
}

// EditReq 编辑用户请求
type EditReq struct {
	UserId     uint   `json:"user_id"`
	NickName   string `json:"nickname"`
	Email      string `json:"email"`
	UserStatus int    `json:"user_status"` // 用户状态 0全部 1正常 2封禁
	RoleIds    []uint `json:"role_ids"`
	AvatarPath string `json:"avatar_path"`
}

// DeleteReq 删除用户请求
type DeleteReq struct {
	Ids []uint `json:"ids"`
}

// ListReq 获取所有用户列表请求
type ListReq struct {
	Page               int    `json:"page"`
	Limit              int    `json:"limit"`
	NickName           string `json:"nickname"`
	Email              string `json:"email"`
	UserStatus         int    `json:"user_status"`
	Heat               int    `json:"heat"`
	FansCount          int    `json:"fans_count"`
	RoleIds            []uint `json:"role_ids"`
	LastLoginTimeBegin string `json:"last_login_time_begin"`
	LastLoginTimeEnd   string `json:"last_login_time_end"`
	CreateTimeBegin    string `json:"create_time_begin"`
	CreateTimeEnd      string `json:"create_time_end"`
}

// ListRes 获取所有用户列表响应
type ListRes struct {
	Id            uint
	AvatarPath    string `json:"avatar_path"` // 头像路径
	NickName      string `json:"nickname"`
	Email         string `json:"email"`
	Heat          int    `json:"heat"`
	FansCount     int    `json:"fans_count"`
	RoleIds       []uint `json:"role_ids"`
	UserStatus    int    `json:"user_status"` // 用户状态 1正常 2封禁 0全部
	LastLoginTime string `json:"last_login_time"`
	CreateTime    string `json:"create_time"`
}

// GetInfoRes 获取当前用户基本信息响应
type GetInfoRes struct {
	Id         uint
	AvatarPath string `json:"avatar_path"` // 头像路径
	NickName   string `json:"nickname"`
	Email      string `json:"email"`
	UserStatus int    `json:"user_status"` // 用户状态 1正常 2封禁 0全部
	RoleIds    []uint `json:"role_ids"`
}
