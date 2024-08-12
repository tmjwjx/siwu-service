package repositories

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"gorm.io/gorm"
	"reflect"
)

// 注意：最好不要使用 .Update方法，Update 方法只更新指定的列，不会自动更新 UpdatedAt 字段。可以使用 .Save 或 .Updates 方法来自动更新 UpdatedAt 字段。

// // InsertUser 创建用户
// func InsertUser(db *gorm.DB, user *models.User) error {
// 	// 向数据库中插入新的用户
// 	result := db.Create(user)
// 	// 检查插入是否成功
// 	if result.Error != nil {
// 		return fmt.Errorf("InsertUser err: 插入用户失败")
// 	}
// 	return nil
// }

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

// InsertFollow 向中间表中插入关注与被关注信息
func InsertFollow(db *gorm.DB, followerId uint, followedId uint) error {
	// 构建要插入的数据
	followData := map[string]interface{}{
		"follower_id": followerId,
		"followed_id": followedId,
	}

	// 插入数据
	result := db.Table("user_follows").Create(followData)

	// 检查是否有错误
	if result.Error != nil {
		return fmt.Errorf("InsertFollow() err: %v", result.Error)
	}
	// 检查是否插入成功
	if result.RowsAffected == 0 {
		return fmt.Errorf("InsertFollow() err: 插入关注记录失败")
	}

	return nil
}

// DeleteObjects 根据一个或多个条件删除一个或多个对象。
// db: GORM 的数据库实例。
// modelType: 要删除的模型类型的指针。
// condition: 删除条件的键值对。
// int64: 返回删除的记录数。
func DeleteObjects(db *gorm.DB, modelType interface{}, condition map[string]interface{}) (int64, error) {
	// 确保 modelType 是指针类型
	if reflect.TypeOf(modelType).Kind() != reflect.Ptr {
		return 0, fmt.Errorf("DeleteObjects() err: 数据模型必须是指针类型")
	}

	// 构建查询条件
	query := db.Model(modelType)
	for key, value := range condition {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// 执行删除操作
	result := query.Delete(modelType)
	if result.Error != nil {
		return 0, fmt.Errorf("DeleteObjects() err: %v\t执行的查询语句为: %v", result.Error, result.Statement.SQL.String())
	}

	// 返回删除的记录数
	return result.RowsAffected, nil
}

// UpdateObjects 根据一个或多个参数更新一个或多个对象。
// db: GORM 的数据库实例。
// modelType: 要更新的模型类型的指针。
// condition: 查询条件的键值对。
// updates: 更新值的键值对。
func UpdateObjects(db *gorm.DB, modelType interface{}, condition map[string]interface{}, updates map[string]interface{}) error {
	// 确保模型类型是指针类型
	modelValue := reflect.ValueOf(modelType)
	if modelValue.Kind() != reflect.Ptr {
		return fmt.Errorf("UpdateObjects() err: 数据模型必须是指针类型")
	}

	// 获取模型的实例
	modelInstance := reflect.New(modelValue.Elem().Type()).Interface()

	// 构建查询条件
	query := db.Model(modelInstance)
	for key, value := range condition {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// 执行更新
	result := query.Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("UpdateObjects() err: %v", result.Error)
	}

	return nil
}

// QueryUserByEmail 通过email查找用户
func QueryUserByEmail(db *gorm.DB, email string) *models.User {
	var user models.User
	d := db.Model(&models.User{}).Where("email = ?", email).Select("*").Scan(&user)
	if d.RowsAffected == 0 {
		return nil
	}
	return &user
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

// // QueryFollowed 查询某id被谁关注了被关注者
// func QueryFollowed(db *gorm.DB, follower uint) {
//  }
// // QueryFollower 查询关注者
// func QueryFollower(db *gorm.DB, followed uint) {
//
// }

// QueryLastUserVerifyCodeByUserID 根据 UserID 查询最后一条 UserVerifyCode 记录（不论该数据的DeleteAt是否已经被赋值）
func QueryLastUserVerifyCodeByUserID(db *gorm.DB, userID uint) (*models.UserVerifyCode, error) {
	var userVerifyCode models.UserVerifyCode
	// 查询最后一条创建的记录（Unscoped()：不论这条数据的deleteAt是否被赋值）
	result := db.Unscoped().Where("user_id = ?", userID).Order("created_at DESC").First(&userVerifyCode)
	if result.Error != nil {
		return nil, fmt.Errorf("QueryLastUserVerifyCodeByUserID() err: 查找id为%d用户验证码失败, 执行的查询语句为: %v", userID, result.Statement.SQL.String())
	}

	return &userVerifyCode, nil
}

// QueryObject 根据一个或多个条件查询一个对象。
// db: GORM 的数据库实例。
// modelType: 要更新的模型类型的指针。
// condition: 查询条件的键值对。
// interface{}: 返回指针类型。
func QueryObject(db *gorm.DB, modelType interface{}, condition map[string]interface{}) (interface{}, error) {
	// 确保模型类型是指针类型
	modelValue := reflect.ValueOf(modelType)
	if modelValue.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("QueryObject() err: 数据模型必须是指针类型")
	}

	// 创建一个新实例，用于存储查询结果
	result := reflect.New(reflect.TypeOf(modelType).Elem()).Interface()

	// 构建查询
	query := db.Model(modelType)
	for key, value := range condition {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	queryResult := query.First(result)
	// 如果错误是“record not found”
	if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
		fmt.Println("QueryObject() 记录未找到")
		return nil, nil
	} else if queryResult.Error != nil { // 如果是其他错误
		return nil, fmt.Errorf("QueryObject() err: 查找对象失败\t执行的查询语句为: %v", queryResult.Statement.SQL.String())
	}

	return result, nil
}

// QueryObjects 根据一个或多个条件查询一个或多个对象。
// db: GORM 的数据库实例。
// modelType: 不是指针类型。
// condition: 查询条件的键值对。
// interface{}: 返回结构体实例切片。
func QueryObjects(db *gorm.DB, modelType interface{}, condition map[string]interface{}) (interface{}, error) {
	// 获取模型的类型，并创建切片类型的实例
	modelTypeValue := reflect.TypeOf(modelType)
	sliceType := reflect.SliceOf(modelTypeValue)
	sliceValue := reflect.New(sliceType).Interface()

	// 构建查询条件
	query := db.Model(modelType)
	for key, value := range condition {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// 执行查询，将结果存储到切片中
	result := query.Find(sliceValue)
	if result.Error != nil {
		return nil, fmt.Errorf("QueryObjects() err: %v\t执行的查询语句为: %v", result.Error, query.Statement.SQL.String())
	}

	return sliceValue, nil
}
