package repositories

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/utils"
	"gorm.io/gorm"
)

// QueryRolesByPage 检索角色，支持分页
// conditions: 查询条件。
// page: 第几页。
// limit: 每页数据条数。
// 例子：如果 page = 2，limit = 10，那么会跳过前 10 条记录，返回第 11-20 条记录。
// func QueryRolesByPage(db *gorm.DB, conditions map[string]interface{}, page int, limit int) ([]*models.Role, error) {

// UpdateAdminRoles 更新用户的角色，确保 userId 存在
func UpdateAdminRoles(db *gorm.DB, userId uint, newRoleIds []uint) error {
	var currentRoles []models.AdminRole
	var currentRoleIds []uint // 当前 userId 所拥有的角色id

	// 查询当前 userId 对应的 RoleIds
	err := db.Where("user_id = ?", userId).Find(&currentRoles).Error
	if err != nil {
		return fmt.Errorf("UpdateAdminRoles() err: %v", err)
	}

	// 该用户没有角色
	if len(currentRoles) == 0 {

	} else {
		// 提取当前的 RoleIds
		for _, role := range currentRoles {
			currentRoleIds = append(currentRoleIds, role.RoleId)
		}

		var deleteRoleIds []uint // 要删除的role_id
		var addRoleIds []uint    // 要添加的role_id

		// 找出需要删除的 RoleIds (当前有的，但不在新列表中)
		for _, v := range currentRoleIds {
			if !utils.IsUintSliContainUint(newRoleIds, v) {
				deleteRoleIds = append(deleteRoleIds, v)
			}
		}

		// 找出需要新增的 RoleIds (新列表中有的，但当前没有)
		for _, v := range newRoleIds {
			if !utils.IsUintSliContainUint(currentRoleIds, v) {
				addRoleIds = append(addRoleIds, v)
			}
		}

		// 删除不在新列表中的角色
		for _, v := range deleteRoleIds {
			_, err := DeleteObjectsByModel(db, &models.AdminRole{}, map[string]interface{}{"user_id": userId, "role_id": v})
			if err != nil {
				return fmt.Errorf("UpdateAdminRoles() err: %v", err)
			}
		}

		// 插入新角色
		for _, v := range addRoleIds {
			err = InsertObject(db, &models.AdminRole{AdminId: userId, RoleId: v})
			if err != nil {
				return fmt.Errorf("UpdateAdminRoles() err: %v", err)
			}
		}
	}
	return nil
}

// QueryRoleById 通过角色id查询该角色的信息
func QueryRoleById(db *gorm.DB, id uint) *models.Role {
	var role models.Role
	d := db.Model(&models.Role{}).Where("id = ?", id).Select("*").Scan(&role)
	if d.RowsAffected == 0 {
		return nil
	}
	return &role
}
