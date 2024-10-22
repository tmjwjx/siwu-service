package repositories

import (
	"forum/internal/models"
	"gorm.io/gorm"
)

// QueryUserById
// @Description: 通过ID查找用户。
// @Author lizhuang 2024-10-21 20:32:02
// @param        db *gorm.DB
// @param        id uint
// @return       *models.User
func QueryUserById(db *gorm.DB, id uint) *models.User {
	var user models.User
	d := db.Model(&models.User{}).Where("id = ?", id).Select("*").Scan(&user)
	// 没有找到用户
	if d.RowsAffected <= 0 {
		return nil
	}
	return &user
}

// QueryUserByEmail
// @Description: 通过email查找用户。
// @Author lizhuang 2024-10-21 20:29:50
// @param        db *gorm.DB
// @param        email string
// @return       *models.User
func QueryUserByEmail(db *gorm.DB, email string) *models.User {
	var user models.User
	d := db.Model(&models.User{}).Where("email = ?", email).Select("*").Scan(&user)
	// 没有找到用户
	if d.RowsAffected <= 0 {
		return nil
	}
	return &user
}

// QueryAdminRoleByUserId
// @Description: 查询用户拥有的角色id
// @Author lizhuang 2024-10-21 20:59:19
// @param        db *gorm.DB
// @param        userId uint
// @return       []uint
func QueryAdminRoleByUserId(db *gorm.DB, userId uint) []uint {
	var roleIds []uint
	db.Model(&models.AdminRole{}).Where("admin_id = ?", userId).Select("role_id").Scan(&roleIds)
	return roleIds
}

// QueryRoleById
// @Description: 通过角色id查询该角色的信息
// @Author lizhuang 2024-10-21 20:59:15
// @param        db *gorm.DB
// @param        id uint
// @return       *models.Role
func QueryRoleById(db *gorm.DB, id uint) *models.Role {
	var role models.Role
	d := db.Model(&models.Role{}).Where("id = ?", id).Select("*").Scan(&role)
	if d.RowsAffected <= 0 {
		return nil
	}
	return &role
}
