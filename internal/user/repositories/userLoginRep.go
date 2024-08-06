package repositories

import (
	"errors"
	"forum/internal/models"
	"gorm.io/gorm"
	"reflect"
)

// JudgeEmailExist 判断某个邮箱是否存在。返回true：存在；false：不存在
func JudgeEmailExist(db *gorm.DB, email string) bool {
	d := db.Table("users").Where("email = ?", email).Select("id")
	if d.RowsAffected != 0 {
		return true
	}
	return false
}

// Create 插入新数据。data应该是指针类型。
func Create(db *gorm.DB, data interface{}) error {
	// 检查 data 是否为指针类型
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		return errors.New("数据参数必须是指针类型")
	}

	// 如果是指针类型，则插入数据
	result := db.Create(data)
	// 检查插入是否成功
	if result.Error != nil {
		return errors.New("插入新数据")
	}
	return nil
}

// CreateUser 创建用户
func CreateUser(db *gorm.DB, user *models.User) error {
	// 向数据库中插入新的用户
	result := db.Create(user)
	// 检查插入是否成功
	if result.Error != nil {
		return errors.New("插入用户失败")
	}
	return nil
}
