package repositories

import (
	"fmt"
	"forum/internal/models"
	"gorm.io/gorm"
)

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

// QueryAdminByEmail 通过email查找管理者
func QueryAdminByEmail(db *gorm.DB, email string) *models.Administrator {
	var admin models.Administrator
	d := db.Model(&models.Administrator{}).Where("email = ?", email).Select("*").Scan(&admin)
	// 没有找到用户
	if d.RowsAffected <= 0 {
		return nil
	}
	return &admin
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

// QuerySuperAdminId
// @Description: 查询超级管理员的ID
// @Author wangyulong 2025-02-07 17:55:04
// @param        db *gorm.DB
// @return       uint
// @return       error
func QuerySuperAdminId(db *gorm.DB) (uint, error) {
	var superAdminId uint

	// 查询超级管理员的id
	err := db.Model(&models.Role{}).Select("id").Where("name = ?", "超级管理员").Scan(&superAdminId).Error
	if err != nil {
		return 0, fmt.Errorf("GetMenuPermRep -> 查询超级管理员的id异常 -> %s", err)
	}

	if superAdminId == 0 {
		return 0, fmt.Errorf("GetMenuPermRep -> 不存在超级管理员这个角色 -> %s", err)
	}

	return superAdminId, nil
}
