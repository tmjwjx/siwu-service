package repositories

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/utils"
	"gorm.io/gorm"
)

// UpdateAdminRoles 更新用户的角色（要确保 userId 存在）
func UpdateAdminRoles(db *gorm.DB, userId uint, newRoleIds []uint) error {
	// 查询当前 userId 拥有的 RoleIds
	var currentRoleIds []uint
	err := db.Model(&models.AdminRole{}).Select("role_id").Where("admin_id = ?", userId).Find(&currentRoleIds).Error
	if err != nil {
		return fmt.Errorf("UpdateAdminRoles() err: %v", err)
	}

	var deleteRoleIds []uint // 要删除的role_id
	var addRoleIds []uint    // 要添加的role_id

	// 找出 userId 需要删除的 RoleIds (当前有的，但不在新列表中)
	for _, v := range currentRoleIds {
		if !utils.IsUintSliContainUint(newRoleIds, v) {
			deleteRoleIds = append(deleteRoleIds, v)
		}
	}
	// 找出 userId 需要增加的 RoleIds (新列表中有的，但当前没有)
	for _, v := range newRoleIds {
		if !utils.IsUintSliContainUint(currentRoleIds, v) {
			addRoleIds = append(addRoleIds, v)
		}
	}

	// 删除
	for _, v := range deleteRoleIds {
		_, err := DeleteObjectsByModel(db, &models.AdminRole{}, map[string]interface{}{"admin_id": userId, "role_id": v})
		if err != nil {
			return fmt.Errorf("UpdateAdminRoles() err: %v", err)
		}
	}
	// 插入
	for _, v := range addRoleIds {
		err = InsertObject(db, &models.AdminRole{AdminId: userId, RoleId: v})
		if err != nil {
			return fmt.Errorf("UpdateAdminRoles() err: %v", err)
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

// QueryUserListByPage 分页查询所有用户列表
// conditions: 查询条件。
// page: 第几页。
// limit: 每页数据条数。
// 例子：如果 page = 2，limit = 10，那么会跳过前 10 条记录，返回第 11-20 条记录。
func QueryUserListByPage(db *gorm.DB, conditions map[string]interface{}, page int, limit int) ([]*models.User, error) {
	var users []*models.User

	// 使用条件查询
	query := db.Model(&models.User{})
	for key, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// 添加分页逻辑
	if page > 0 && limit > 0 {
		offset := (page - 1) * limit // 计算偏移量。Offset 是从第几条数据开始取。
		query = query.Offset(offset).Limit(limit)
	}

	// 执行查询
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
