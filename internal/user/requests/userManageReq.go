package requests

type ReseatReq struct {
	Id uint `json:"id"`
}

// AddAndEditReq 添加、编辑用户请求
type AddAndEditReq struct {
	NickName   string `json:"nick_name"`
	Email      string `json:"email"`
	UserStatus int    `json:"user_status"` // 用户状态 1正常 2封禁 0全部
	RoleIds    []uint `json:"role_ids"`
}

// DeleteReq 删除用户请求
type DeleteReq struct {
	Ids []uint `json:"ids"`
}

// ListReq 获取所有用户列表请求
type ListReq struct {
	Page          int    `json:"page"`
	Limit         int    `json:"limit"`
	NickName      string `json:"nick_name"`
	Email         string `json:"email"`
	UserStatus    int    `json:"user_status"`
	Heat          int    `json:"heat"`
	FansCount     int    `json:"fans_count"`
	RoleIds       []uint `json:"role_ids"`
	LastLoginTime string `json:"last_login_time"`
	CreateTime    string `json:"create_time"`
}

// ListRes 获取所有用户列表响应
type ListRes struct {
	Id            uint
	AvatarPath    string `json:"avatar_path"` // 头像路径
	NickName      string `json:"nick_name"`
	Email         string `json:"email"`
	Heat          int    `json:"heat"`
	FansCount     int    `json:"fans_count"`
	Roles         []Role `json:"roles"`
	UserStatus    int    `json:"user_status"` // 用户状态 1正常 2封禁 0全部
	LastLoginTime string `json:"last_login_time"`
	CreateTime    string `json:"create_time"`
}

// Role 角色
type Role struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

// GetInfoRes 获取当前用户基本信息响应
type GetInfoRes struct {
	Id            uint
	AvatarPath    string `json:"avatar_path"` // 头像路径
	NickName      string `json:"nick_name"`
	Email         string `json:"email"`
	UserStatus    int    `json:"user_status"` // 用户状态 1正常 2封禁 0全部
	Roles         []Role `json:"roles"`
	LastLoginTime string `json:"last_login_time"`
	CreateTime    string `json:"create_time"`
}
