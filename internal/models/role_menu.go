package models

// RoleMenu
// @Description: 角色菜单关联表
// @Author wangyulong 2024-10-19 16:20:08
type RoleMenu struct {
	RoleId uint `json:"role_id"`
	MenuId uint `json:"menu_id"`
}
