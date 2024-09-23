package repositories

import (
	"fmt"
	"forum/internal/models"
	"gorm.io/gorm"
	"reflect"
)

// InsertObject 插入一条数据。
// db: GORM 的数据库实例。
// model: 指针类型。
func InsertObject(db *gorm.DB, model interface{}) error {
	// 检查 data 是否为指针类型
	if reflect.TypeOf(model).Kind() != reflect.Ptr {
		return fmt.Errorf("InsertObject() err: 数据模型必须是指针类型")
	}

	// 插入数据
	result := db.Create(model)
	if result.Error != nil {
		return fmt.Errorf("InsertObject() err: %v\t执行的查询语句为: %v", result.Error, result.Statement.SQL.String())
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("InsertObject() err: 插入数据失败")
	}

	return nil
}

// InsertObjects 插入多条数据。
// db: GORM 的数据库实例。
// models: 切片类型。
func InsertObjects(db *gorm.DB, models interface{}) error {
	// 确保传入的 models 是一个切片
	modelsValue := reflect.ValueOf(models)
	if modelsValue.Kind() != reflect.Slice {
		return fmt.Errorf("InsertObjects() err: 传入的参数必须是切片类型")
	}

	// 执行批量插入
	result := db.Create(models)
	if result.Error != nil {
		return fmt.Errorf("InsertObjects() err: %v\t执行的查询语句为: %v", result.Error, result.Statement.SQL.String())
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("InsertObjects() err: 插入数据失败")
	}
	return nil
}

// DeleteObjectsByModel 按照model模型，根据一个或多个条件删除一个或多个对象。
// db: GORM 的数据库实例。
// modelType: 要删除的模型类型的指针。
// condition: 删除条件的键值对。
// int64: 返回删除的记录数。
func DeleteObjectsByModel(db *gorm.DB, modelType interface{}, condition map[string]interface{}) (int64, error) {
	// 确保 modelType 是指针类型
	if reflect.TypeOf(modelType).Kind() != reflect.Ptr {
		return 0, fmt.Errorf("DeleteObjectsByModel() err: 数据模型必须是指针类型")
	}

	// 构建查询条件
	query := db.Model(modelType)
	for key, value := range condition {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// 执行删除操作
	result := query.Delete(modelType)
	if result.Error != nil {
		return 0, fmt.Errorf("DeleteObjectsByModel() err: %v\t执行的查询语句为: %v", result.Error, result.Statement.SQL.String())
	}

	// 返回删除的记录数
	return result.RowsAffected, nil
}

// UpdateObjects 根据一个或多个参数更新一个或多个对象。
// db: GORM 的数据库实例。
// model: 要更新的模型类型的指针，模型中要包含查询的参数。
// updates: 更新值的键值对。
func UpdateObjects(db *gorm.DB, model interface{}, updates map[string]interface{}) error {
	// 检查 model 是否为指针类型
	if reflect.TypeOf(model).Kind() != reflect.Ptr {
		return fmt.Errorf("UpdateObjects() err: 数据模型必须是指针类型")
	}

	// 执行更新操作
	result := db.Model(model).Where(model).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("UpdateObjects() err: %v", result.Error)
	}

	return nil
}

// QueryRoles 检索角色。
// conditions: 查询条件。
func QueryRoles(db *gorm.DB, conditions map[string]interface{}) ([]*models.Role, error) {
	var roles []*models.Role

	// 使用条件查询
	query := db.Model(&models.Role{})
	for key, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// 执行查询
	if err := query.Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

// QueryRolesByPage 检索角色，支持分页
// conditions: 查询条件。
// page: 第几页。
// limit: 每页数据条数。
// 例子：如果 page = 2，limit = 10，那么会跳过前 10 条记录，返回第 11-20 条记录。
// 返回的int表示一共有多少条符合条件的数据
func QueryRolesByPage(db *gorm.DB, conditions map[string]interface{}, page int, limit int) ([]*models.Role, int, error) {
	var roles []*models.Role

	// 使用条件查询
	query := db.Model(&models.Role{})
	for key, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// 查看符合条件的数据一共有多少条
	query.Find(&roles)
	total := len(roles)

	// 添加分页逻辑
	if page > 0 && limit > 0 {
		offset := (page - 1) * limit // 计算偏移量。Offset 是从第几条数据开始取。
		query = query.Offset(offset).Limit(limit)
	}

	// 执行查询
	if err := query.Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// QueryRoleById 通过ID查找角色
func QueryRoleById(db *gorm.DB, id uint) *models.Role {
	var role models.Role
	d := db.Model(&models.Role{}).Where("id = ?", id).Select("*").Scan(&role)
	if d.RowsAffected == 0 {
		return nil
	}

	return &role
}

// QueryUserById 通过ID查找用户
func QueryUserById(db *gorm.DB, id uint) *models.User {
	var user models.User
	d := db.Model(&models.User{}).Where("id = ?", id).Select("*").Scan(&user)
	if d.RowsAffected == 0 {
		return nil
	}
	return &user
}
