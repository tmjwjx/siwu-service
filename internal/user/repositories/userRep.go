package repositories

import (
	"fmt"
	"forum/internal/models"
	"gorm.io/gorm"
	"reflect"
)

// 用户共用方法

// InsertObject 插入一条数据。
//
// db: GORM 的数据库实例。
// model: 指针类型。
//
// Returns an error
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
//
// db: GORM 的数据库实例。
// models: 切片类型。
//
// Returns an error
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
//
// db: GORM 的数据库实例。
// modelType: 要删除的模型类型的指针。
// condition: 删除条件的键值对。
//
// Returns int64:返回删除的记录数。
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

// DeleteObjectsByTable 按照表名，根据一个或多个条件删除一个或多个对象。
//
// db: GORM 的数据库实例。
// tableName: 要删除的表名。
// condition: 删除条件的键值对。
//
// Returns int64: 返回删除的记录数；error: 错误
func DeleteObjectsByTable(db *gorm.DB, tableName string, condition map[string]interface{}) (int64, error) {
	// 构建查询条件
	query := db.Table(tableName)
	for key, value := range condition {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// 执行删除操作。query.Delete(nil)：删除符合条件的记录，不需要指定具体的模型类型。
	result := query.Delete(nil)
	if result.Error != nil {
		return 0, fmt.Errorf("DeleteObjectsByModel() err: %v\t执行的查询语句为: %v", result.Error, result.Statement.SQL.String())
	}

	// 返回删除的记录数
	return result.RowsAffected, nil
}

// UpdateObjects 根据一个或多个参数更新一个或多个对象。
//
// db: GORM 的数据库实例。
// model: 要更新的模型类型的指针，模型中要包含查询的参数。
// updates: 更新值的键值对。
//
// Returns error: 错误
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

// QueryUserById 通过ID查找用户
func QueryUserById(db *gorm.DB, id uint) *models.User {
	var user models.User
	d := db.Model(&models.User{}).Where("id = ?", id).Select("*").Scan(&user)
	// 没有找到用户
	if d.RowsAffected <= 0 {
		return nil
	}
	return &user
}

// QueryUserByEmail 通过email查找用户
func QueryUserByEmail(db *gorm.DB, email string) *models.User {
	var user models.User
	d := db.Model(&models.User{}).Where("email = ?", email).Select("*").Scan(&user)
	// 没有找到用户
	if d.RowsAffected <= 0 {
		return nil
	}
	return &user
}

// QueryUserByIdIncludeSoftDelete 通过ID查找用户，包括已软删除的用户）
func QueryUserByIdIncludeSoftDelete(db *gorm.DB, id uint) *models.User {
	var user models.User
	// 使用 Unscoped() 包括软删除的记录
	d := db.Unscoped().Model(&models.User{}).Where("id = ?", id).Select("*").Scan(&user)

	// 没有找到用户
	if d.RowsAffected <= 0 {
		return nil
	}
	return &user
}

// QueryUserDetailsById 通过用户ID查找用户详情
func QueryUserDetailsById(db *gorm.DB, id uint) *models.UserDetail {
	var userDetail models.UserDetail
	d := db.Model(&models.UserDetail{}).Where("user_id = ?", id).Select("*").Scan(&userDetail)
	if d.RowsAffected <= 0 {
		return nil
	}
	return &userDetail
}

// QueryFollowed 查询id关注了谁。（查看我的关注）
func QueryFollowed(db *gorm.DB, follower uint) ([]uint, error) {
	var followedIDSli []uint
	// 执行查询，获取所有关注的用户ID
	// Pluck("followed_id", &followedIDs)：只提取 followed 字段的值，并存储到 followedIDs 切片中。
	err := db.Table("sw_user_follows").Where("follower_id = ?", follower).Pluck("followed_id", &followedIDSli).Error
	if err != nil {
		return nil, fmt.Errorf("QueryFollowed() err: %v", err)
	}

	return followedIDSli, nil
}

// QueryFollower 查询id被谁关注。（查看我的粉丝）
func QueryFollower(db *gorm.DB, followed uint) ([]uint, error) {
	var followerIDSli []uint
	// 执行查询，获取所有关注的用户ID
	err := db.Table("sw_user_follows").Where("followed_id = ?", followed).Pluck("follower_id", &followerIDSli).Error
	if err != nil {
		return nil, fmt.Errorf("QueryFollowed() err: %v", err)
	}

	return followerIDSli, nil
}
