package repositories

// // QueryLastLoginTime 查询用户的最后登录时间
// func QueryLastLoginTime(db *gorm.DB, userID uint) (time.Time, error) {
// 	var user models.User
//
// 	// 查询用户的最后登录时间
// 	if err := db.Model(&models.User{}).Select("last_login_time").Where("id = ?", userID).First(&user).Error; err != nil {
// 		return time.Time{}, fmt.Errorf("QueryLastLoginTime() err: %v", err)
// 	}
//
// 	return user.LastLoginTime, nil
// }
