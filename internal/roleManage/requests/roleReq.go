package requests

// RoleReq 添加角色请求
type RoleReq struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Status int    `json:"status"`
	Sort   int    `json:"sort"`
}

// RoleIdsReq 根据Id切片删除Role请求
type RoleIdsReq struct {
	Ids []uint `json:"ids"`
}

// SearchRoleReq 检索Role请求
type SearchRoleReq struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Status int    `json:"status"`
	// Sort   int    `json:"sort"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// SearchRoleRes 检索角色响应
type SearchRoleRes struct {
	Id        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Status    int    `json:"status"`
	Sort      int    `json:"sort"`
}

// DispatchRoleReq 为用户分配角色请求
type DispatchRoleReq struct {
	UserId uint   `json:"user_id"`
	Ids    []uint `json:"ids"`
}
