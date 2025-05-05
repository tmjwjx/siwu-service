package repositories

import (
	"forum/internal/models"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

// QueryUserById
// @Description: 通过ID查找用户。
// @Author lizhuang 2024-10-04 21:27:16
// @param        db *gorm.DB
// @param        id uint
// @return       *models.User
func QueryUserById(db *gorm.DB, id uint) *models.User {
	var user models.User
	d := db.Model(&models.User{}).Where("id = ?", id).Select("*").Scan(&user)
	// 没有找到用户
	if d.RowsAffected <= 0 {
		return nil
	}
	return &user
}

// QueryUserByEmail
// @Description: 通过email查找用户。
// @Author lizhuang 2024-10-04 21:27:26
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

// QueryUserByIdIncludeSoftDelete
// @Description: 通过ID查找用户，包括已软删除的用户。
// @Author lizhuang 2024-10-04 21:27:33
// @param        db *gorm.DB
// @param        id uint
// @return       *models.User
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

// QueryUserDetailsById
// @Description: 通过用户ID查找用户详情。
// @Author lizhuang 2024-10-04 21:27:42
// @param        db *gorm.DB
// @param        id uint
// @return       *models.UserDetail
func QueryUserDetailsById(db *gorm.DB, id uint) *models.UserDetail {
	var userDetail models.UserDetail
	d := db.Model(&models.UserDetail{}).Where("user_id = ?", id).Select("*").Scan(&userDetail)
	if d.RowsAffected <= 0 {
		return nil
	}
	return &userDetail
}

// QueryFollowed
// @Description: 查询id关注了谁。
// @Author lizhuang 2024-10-04 21:27:56
// @param        db *gorm.DB
// @param        follower uint
// @return       []uint
// @return       error
func QueryFollowed(db *gorm.DB, follower uint) ([]uint, error) {
	var followedIDSli []uint
	// 执行查询，获取所有关注的用户ID
	// Pluck("followed_id", &followedIDs)：只提取 followed 字段的值，并存储到 followedIDs 切片中。
	err := db.Table("sw_user_follows").Where("follower_id = ?", follower).Pluck("followed_id", &followedIDSli).Error
	if err != nil {
		// return nil, fmt.Errorf("QueryFollowed() err: %v", err)
		globals.Log.Error(err.Error())
		return nil, err
	}

	return followedIDSli, nil
}

// QueryFollower
// @Description: 查询id被谁关注。
// @Author lizhuang 2024-10-04 21:28:04
// @param        db *gorm.DB
// @param        followed uint
// @return       []uint
// @return       error
func QueryFollower(db *gorm.DB, followed uint) ([]uint, error) {
	var followerIDSli []uint
	// 执行查询，获取所有关注的用户ID
	err := db.Table("sw_user_follows").Where("followed_id = ?", followed).Pluck("follower_id", &followerIDSli).Error
	if err != nil {
		// return nil, fmt.Errorf("QueryFollowed() err: %v", err)
		globals.Log.Error(err.Error())
		return nil, err
	}

	return followerIDSli, nil
}
