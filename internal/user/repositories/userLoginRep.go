package repositories

import (
	"fmt"
	"forum/internal/models"
	"gorm.io/gorm"
	"reflect"
)

// QueryUserByEmail 通过email查找用户
func QueryUserByEmail(db *gorm.DB, email string) *models.User {
	var user *models.User
	d := db.Table("t_users").Where("email = ?", email).Select("*").Scan(user)
	if d.RowsAffected == 0 {
		return nil
	}
	return user
}

// QueryUserVerifyCodeByUID 根据UserID查询验证码信息
func QueryUserVerifyCodeByUID(db *gorm.DB, userId uint) *models.UserVerifyCode {
	var userVerifyCode *models.UserVerifyCode
	d := db.Table("t_user_verify_code").Where("user_id = ?", userId).Select("*").Scan(userVerifyCode)
	if d.RowsAffected == 0 {
		return nil
	}
	return userVerifyCode
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
		return fmt.Errorf("Insert err: 插入新数据失败")
	}
	return nil
}

// InsertUser 创建用户
func InsertUser(db *gorm.DB, user *models.User) error {
	// 向数据库中插入新的用户
	result := db.Create(user)
	// 检查插入是否成功
	if result.Error != nil {
		return fmt.Errorf("InsertUser err: 插入用户失败")
	}
	return nil
}

// UpdateVerifyCodeByUID 根据userId更新验证码
func UpdateVerifyCodeByUID(db *gorm.DB, userId uint, verifyCode string) error {
	d := db.Table("t_user_verify_code").Where("user_id = ?", userId).Update("verify_code = ", verifyCode)
	if d.RowsAffected == 0 {
		return fmt.Errorf("UpdateVerifyCodeByUID err: 更新验证码失败")
	}
	return nil
}

// 根据某条件查询全部
