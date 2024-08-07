package repositories

import (
	"fmt"
	"forum/internal/models"
	"gorm.io/gorm"
	"reflect"
)

// 注意：最好不要使用 .Update方法，Update 方法只更新指定的列，不会自动更新 UpdatedAt 字段。可以使用 .Save 或 .Updates 方法来自动更新 UpdatedAt 字段。

// QueryUserByEmail 通过email查找用户
func QueryUserByEmail(db *gorm.DB, email string) *models.User {
	var user models.User
	d := db.Model(&models.User{}).Where("email = ?", email).Select("*").Scan(&user)
	if d.RowsAffected == 0 {
		return nil
	}
	return &user
}

// Insert 插入新数据。data应该是指针类型。
func Insert(db *gorm.DB, data interface{}) error {
	// 检查 data 是否为指针类型
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		return fmt.Errorf("Insert err: 数据参数必须是指针类型")
	}

	// 如果是指针类型，则插入数据
	result := db.Create(data)
	// 检查插入是否成功
	if result.Error != nil {
		return fmt.Errorf("Insert err: 插入新数据失败\t执行的查询语句为: %s", result.Statement.SQL.String())
	}
	return nil
}

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

// QueryUserVerifyCodeByUID 根据UserID查询验证码信息
func QueryUserVerifyCodeByUID(db *gorm.DB, userId uint) *models.UserVerifyCode {
	var userVerifyCode models.UserVerifyCode
	d := db.Model(&models.UserVerifyCode{}).Where("user_id = ?", userId).Select("*").Scan(&userVerifyCode)
	if d.RowsAffected == 0 {
		return nil
	}

	return &userVerifyCode
}

// UpdateVerifyCodeByUID 根据userId更新验证码
func UpdateVerifyCodeByUID(db *gorm.DB, userId uint, verifyCode string) error {
	d := db.Model(&models.UserVerifyCode{}).Where("user_id = ?", userId).Updates(models.UserVerifyCode{
		VerifyCode: verifyCode,
	})
	if d.Error != nil {
		return fmt.Errorf("UpdateVerifyCodeByUID err: 更新验证码失败, %v\t执行的查询语句为: %s", d.Error, d.Statement.SQL.String())
	}
	if d.RowsAffected == 0 {
		return fmt.Errorf("UpdateVerifyCodeByUID err: 没有找到匹配的记录\t执行的查询语句为: %s", d.Statement.SQL.String())
	}

	return nil
}

// Query 根据参数查询整个对象。查询参数可以是多对key-value。返回全部符合的数据。
func Query(db *gorm.DB, modelType interface{}, param map[string]interface{}) (interface{}, error) {
	// 获取模型的类型，并创建切片类型的实例
	modelTypeValue := reflect.TypeOf(modelType)
	sliceType := reflect.SliceOf(modelTypeValue)
	sliceValue := reflect.New(sliceType).Interface()

	// 构建查询条件
	query := db.Model(modelType)
	for key, value := range param {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// 执行查询，将结果存储到切片中
	result := query.Find(sliceValue)
	if result.Error != nil {
		return nil, fmt.Errorf("Query() err: %s\t执行的查询语句为: %s", result.Error, query.Statement.SQL.String())
	}

	return sliceValue, nil
}

// Update 根据多个参数更新对象。
// db：GORM 的数据库实例。modelType：要更新的模型类型的指针。filters：查询条件的键值对。updates：更新值的键值对。
func Update(db *gorm.DB, modelType interface{}, filters map[string]interface{}, updates map[string]interface{}) error {
	// 确保模型类型是指针类型
	modelValue := reflect.ValueOf(modelType)
	if modelValue.Kind() != reflect.Ptr {
		return fmt.Errorf("modelType must be a pointer")
	}

	// 获取模型的实例
	modelInstance := reflect.New(modelValue.Elem().Type()).Interface()

	// 构建查询条件
	query := db.Model(modelInstance)
	for key, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// 执行更新
	result := query.Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no records updated")
	}

	return nil
}
