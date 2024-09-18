package models

// AdminRole 管理员的角色或权限等级。（多对多关系）
type AdminRole struct {
	// 使用 unique 标签强制 AdminId 和 RoleId 的组合唯一
	AdminId uint `json:"user_id" gorm:"uniqueIndex:user_have_role_unique"`
	RoleId  uint `json:"role_id" gorm:"uniqueIndex:user_have_role_unique"` // 对应的角色id
}
