package requests

// Role 添加角色
type Role struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Status int    `json:"status"`
	Sort   int    `json:"sort"`
}

// RoleIds 根据Id切片删除Role
type RoleIds struct {
	Ids []uint `json:"ids"`
}

// SearchRole 检索Role
type SearchRole struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Status int    `json:"status"`
	Sort   int    `json:"sort"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

// DispatchRole 为用户分配角色
type DispatchRole struct {
	UserId uint   `json:"user_id"`
	Ids    []uint `json:"ids"`
}
