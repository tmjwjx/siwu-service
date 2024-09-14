package models

// UserRole 用户的角色（多对多关系）
type UserRole struct {
	// 使用 unique 标签强制 UserId 和 RoleId 的组合唯一
	UserId uint `json:"user_id" gorm:"uniqueIndex:user_role_unique"`
	RoleId uint `json:"role_id" gorm:"uniqueIndex:user_role_unique"` // 用户对应的角色
}
