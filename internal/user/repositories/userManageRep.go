package repositories

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/casbin"
	"forum/pkg/globals"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"time"
)

// // UpdateAdminRoles 更新用户的角色（要确保 userId 存在）
// func UpdateAdminRoles(db *gorm.DB, userId uint, newRoleIds []uint) error {
// 	// 查询当前 userId 拥有的 RoleNames
// 	var currentRoleIds []uint
// 	err := db.Model(&models.AdminRole{}).Select("role_id").Where("admin_id = ?", userId).Find(&currentRoleIds).Error
// 	if err != nil {
// 		return fmt.Errorf("UpdateAdminRoles() err: %v", err)
// 	}
//
// 	var deleteRoleIds []uint // 要删除的role_id
// 	var addRoleIds []uint    // 要添加的role_id
//
// 	// 找出 userId 需要删除的 RoleNames (当前有的，但不在新列表中)
// 	for _, v := range currentRoleIds {
// 		if !lo.Contains(newRoleIds, v) {
// 			deleteRoleIds = append(deleteRoleIds, v)
// 		}
// 	}
// 	// 找出 userId 需要增加的 RoleNames (新列表中有的，但当前没有)
// 	for _, v := range newRoleIds {
// 		if !lo.Contains(currentRoleIds, v) {
// 			addRoleIds = append(addRoleIds, v)
// 		}
// 	}
//
// 	// 删除
// 	for _, v := range deleteRoleIds {
// 		_, err := sqlUtils.DeleteObjectsByModel(db, &models.AdminRole{}, map[string]interface{}{"admin_id": userId, "role_id": v})
// 		if err != nil {
// 			return fmt.Errorf("UpdateAdminRoles() err: %v", err)
// 		}
// 	}
// 	// 插入
// 	for _, v := range addRoleIds {
// 		err = sqlUtils.InsertObject(db, &models.AdminRole{AdminId: userId, RoleId: v})
// 		if err != nil {
// 			return fmt.Errorf("UpdateAdminRoles() err: %v", err)
// 		}
// 	}
//
// 	return nil
// }

// QueryRoleById 通过角色id查询该角色的信息
func QueryRoleById(db *gorm.DB, id uint) *models.Role {
	var role models.Role
	d := db.Model(&models.Role{}).Where("id = ?", id).Select("*").Scan(&role)
	if d.RowsAffected <= 0 {
		return nil
	}
	return &role
}

// QueryUserListByPage 分页查询所有用户列表
// conditions: 查询条件。
// page: 第几页。
// limit: 每页数据条数。
// heat: 用户热度的下限。
// fans_count: 粉丝数的下限。
// 例子：如果 page = 2，limit = 10，那么会跳过前 10 条记录，返回第 11-20 条记录。
// 返回的int表示一共有多少条符合条件的数据
func QueryUserListByPage(db *gorm.DB, conditions map[string]interface{}, page int, limit int, roleIds []uint, heat, fansCount int, createTimeBegin, createTimeEnd, lastLoginTimeBegin, lastLoginTimeEnd string) ([]*models.User, int, error) {
	var users []*models.User

	// 使用条件查询
	query := db.Model(&models.User{})
	for key, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}
	// 添加 heat 和 fans_count 的条件
	query = query.Where("heat >= ? AND fans_count >= ?", heat, fansCount)
	// 解析 lastLoginTimeBegin 和 lastLoginTimeEnd 字符串为 time.Time 类型
	if lastLoginTimeBegin != "" {
		parsedLastLoginTimeBegin, err := time.Parse("2006-01-02", lastLoginTimeBegin)
		if err != nil {
			return nil, 0, fmt.Errorf("QueryUserListByPage() err: lastLoginTimeBegin 解析错误: %v", err)
		}
		query = query.Where("last_login_time >= ?", parsedLastLoginTimeBegin)
	}

	if lastLoginTimeEnd != "" {
		parsedLastLoginTimeEnd, err := time.Parse("2006-01-02", lastLoginTimeEnd)
		if err != nil {
			return nil, 0, fmt.Errorf("QueryUserListByPage() err: lastLoginTimeEnd 解析错误: %v", err)
		}
		query = query.Where("last_login_time <= ?", parsedLastLoginTimeEnd)
	}

	// 解析 createTimeBegin 和 createTimeEnd 字符串为 time.Time 类型
	if createTimeBegin != "" {
		parsedCreateTimeBegin, err := time.Parse("2006-01-02", createTimeBegin)
		if err != nil {
			return nil, 0, fmt.Errorf("QueryUserListByPage() err: createTimeBegin 解析错误: %v", err)
		}
		query = query.Where("created_at >= ?", parsedCreateTimeBegin)
	}
	if createTimeEnd != "" {
		parsedCreateTimeEnd, err := time.Parse("2006-01-02", createTimeEnd)
		if err != nil {
			return nil, 0, fmt.Errorf("QueryUserListByPage() err: createTimeEnd 解析错误: %v", err)
		}
		query = query.Where("created_at <= ?", parsedCreateTimeEnd)
	}

	// 执行查询，获取用户列表
	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// 对符合条件的用户进行 RoleId 查询和过滤
	var filteredUsers []*models.User // 最终的结果
	for _, user := range users {
		var userRoleIds []uint

		// 获取该用户对应的全部角色id
		casbinService, err := casbin.NewCasbinService(globals.DB)
		if err != nil {
			return nil, 0, fmt.Errorf("QueryUserListByPage() %v", err)
		}
		userRoleIds, err = casbinService.GetRolesForUser(user.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("QueryUserListByPage() %v", err)
		}

		// 判断 userRoleIds 是否包含传入的 roleIds
		if lo.EveryBy(roleIds, func(roleId uint) bool { return lo.Contains(userRoleIds, roleId) }) {
			filteredUsers = append(filteredUsers, user) // 如果包含，则将该用户加入结果集
		}
	}

	// 最后对结果集进行分页
	offset := (page - 1) * limit
	if len(filteredUsers) > offset {
		end := offset + limit
		if end > len(filteredUsers) {
			end = len(filteredUsers)
		}
		return filteredUsers[offset:end], len(filteredUsers), nil
	}

	return []*models.User{}, len(filteredUsers), nil
}

// // QueryAdminRoleByUserId 查询用户拥有的角色id
// func QueryAdminRoleByUserId(db *gorm.DB, userId uint) []uint {
// 	var roleIds []uint
// 	db.Model(&models.AdminRole{}).Where("admin_id = ?", userId).Select("role_id").Scan(&roleIds)
// 	return roleIds
// }

// QueryAllUser 查询所有用户
func QueryAllUser(db *gorm.DB) ([]models.User, error) {
	var users []models.User

	if err := db.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("QueryAllUser() err: Failed to retrieve users from database")
	}

	return users, nil
}
