package repositories

import (
	"fmt"
	"forum/internal/image/controllers"
	"forum/internal/image/logics"
	"forum/internal/internal_pkg/internal_utils"
	"forum/internal/models"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

// QueryPersonEmail 查询Email
func QueryPersonEmail(userAccountReq *requests.UserAccountReq, db *gorm.DB) error {

	// 查询 Email 是否唯一
	var user models.User
	err := db.Where("email = ?", userAccountReq.Email).First(&user).Error
	if err != nil {
		return fmt.Errorf("QueryPersonEmail -> %s", err)
	}

	return nil

}

// UserDataRequest 更新用户个人资料
func UserDataRequest(userDataReq *requests.UserDataReq, db *gorm.DB) error {

	var user models.User

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("UserDataRequest -> 开启事务失败 -> %s", tx.Error)
	}

	// 查询该用户是否存在
	err := tx.Model(&models.User{}).Where("id = ?", userDataReq.ID).First(&user).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("UserDataRequest -> 用户表中用户不存在 -> %s", err)
	}

	// 更新 User 表中的 Nickname
	err = tx.Model(&user).Update("Nickname", userDataReq.Nickname).Error

	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("UserDataRequest -> 用户表中没有更新任何记录 -> %s", err)
	}

	// 更新 UserDetail 表

	var userDetail models.UserDetail
	// 查询该用户的外键是否存在
	err = tx.Where("user_id", userDataReq.ID).First(&userDetail).Error

	if err != nil {
		userDetail.ID = userDataReq.ID
		userDetail.CareerDirection = userDataReq.CareerDirection
		userDetail.HomePage = userDataReq.HomePage
		userDetail.Signature = userDataReq.Signature
		//return fmt.Errorf("UserDataRequest -> 用户详情表中用户不存在 -> %s", err)
		err = tx.Create(&userDetail).Error
		if err != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("UserDataRequest -> 用户详情表中数据插入失败 -> %s", err)
		}
	} else {
		// 使用 Map 更新特定字段，如果 UserDataReq 结构体字段与数据库字段不一致时
		updates := map[string]interface{}{
			"careerDirection": userDataReq.CareerDirection,
			"homePage":        userDataReq.HomePage, // 注意这里要使用数据库中的列名
			"signature":       userDataReq.Signature,
		}

		// 更新用户详情表中相应的字段
		err = tx.Model(&userDetail).Select("CareerDirection", "HomePage", "Signature").Updates(updates).Error
		if err != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("UserDataRequest -> 用户详情表中没有更新任何记录 -> %s", err)
		}
	}
	// 更新用户标签
	var existingTags []models.Tag
	err = tx.Where("name IN ?", userDataReq.UserTags).Find(&existingTags).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("UserDataRequest -> 更新用户标签失败 -> %s", err)
	}

	// 找到所有传递过来的标签的ID
	var tagIDs []uint
	for _, tag := range existingTags {
		tagIDs = append(tagIDs, tag.ID)
	}

	// 清除旧的用户标签关联
	err = tx.Where("user_id = ?", userDataReq.ID).Delete(&models.UserTag{}).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("UserDataRequest -> 清除旧的用户标签关联失败 -> %s", err)
	}

	// 添加新的用户标签关联
	for _, tagID := range tagIDs {
		err := tx.Create(&models.UserTag{UserID: userDataReq.ID, TagID: tagID}).Error
		if err != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("UserDataRequest -> 添加新的用户标签关联失败 -> %s", err)
		}
	}

	u := &logics.UrlParam{
		UrlPath: userDataReq.Path,
		Home:    globals.User,
		HomeID:  userDataReq.ID,
	}
	err = controllers.StoreUrlCtrl(u)
	if err != nil {
		return fmt.Errorf("UserDataRequest -> 存储图片的相关信息失败 -> %s", err)
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("UserDataRequest -> 提交事务失败 -> %s", err)
	}

	return nil

}

// UserAccountRequest 更新用户账号设置
func UserAccountRequest(userAccountReq *requests.UserAccountReq, db *gorm.DB) error {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("UserAccountRequest -> 开启事务失败 -> %s", tx.Error)
	}

	// 查询该用户是否存在
	var user models.User
	err := tx.Where("id = ?", userAccountReq.ID).First(&user).Error

	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("UserAccountRequest -> 用户表中用户不存在 -> %s", err)
	}

	// 更新 User 表中的 email , password
	err = tx.Model(&user).Updates(map[string]interface{}{
		"email":    userAccountReq.Email,
		"password": userAccountReq.Password,
	}).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("UserAccountRequest -> 更新 User 表中的 email , password -> %s", err)
	}

	// 查询该用户的外键是否存在
	var userDetail models.UserDetail
	err = tx.Where("user_id = ?", userAccountReq.ID).First(&userDetail).Error

	if err != nil {
		//return fmt.Errorf("UserAccountRequest -> 用户详情表中用户不存在 -> %s", err)
		userDetail.UserID = userAccountReq.ID
		userDetail.BlogLink = userAccountReq.BlogLink
		userDetail.WeiboLink = userAccountReq.WeiboLink
		userDetail.GithubLink = userAccountReq.GithubLink

		err := tx.Create(&userDetail).Error
		if err != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("UserAccountRequest -> 更新 UserDetail 表中的 BlogLink , WeiboLink , GithubLink字段失败 -> %s", err)
		}
	} else {
		// 更新 UserDetail 表中的 BlogLink , WeiboLink , GithubLink
		err = tx.Model(&userDetail).Updates(map[string]interface{}{
			"blog_link":   userAccountReq.BlogLink,
			"weibo_link":  userAccountReq.WeiboLink,
			"github_link": userAccountReq.GithubLink,
		}).Error

		if err != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("UserAccountRequest -> 更新 UserDetail 表中的 BlogLink , WeiboLink , GithubLink字段失败 -> %s", err)
		}
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("UserAccountRequest -> 提交事务失败 -> %s", err)
	}

	return nil

}

// UserPrivateSetRequest 更新用户私信设置
func UserPrivateSetRequest(userPrivateSetReq *requests.UserPrivateSettingsReq, db *gorm.DB) error {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("UserPrivateSetRequest -> 开启事务失败 -> %s", tx.Error)
	}

	var user models.User
	// 查询该用户是否存在
	err := tx.Model(&models.User{}).Where("id = ?", userPrivateSetReq.ID).First(&user).Error
	if err != nil {
		//return fmt.Errorf("UserPrivateSetRequest -> %s", err)
		// 数据库表中还没有该用户的数据，直接插入即可
		user.ID = userPrivateSetReq.ID
		user.PrivateSettings = userPrivateSetReq.PrivateSettings
		err := tx.Create(&user).Error
		if err != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("UserPrivateSetRequest -> 用户私信设置插入数据失败1 -> %s", err)
		}
	} else {
		// 数据库表中已经存在该用户的信息，直接更新用户信息
		err := tx.Model(&user).Select("PrivateSettings").Updates(userPrivateSetReq).Error
		if err != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("UserPrivateSetRequest -> 用户私信设置插入数据失败2 -> %s", err)
		}
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("UserPrivateSetRequest -> 提交事务失败 -> %s", err)
	}

	return nil

}

// UserDataResponse 响应用户个人资料
func UserDataResponse(userID uint, db *gorm.DB) (*requests.UserDataRes, error) {
	var user models.User

	err := db.Preload("UserDetail").Preload("Tags").First(&user, userID).Error
	if err != nil {
		return nil, fmt.Errorf("UserDataResponse -> %s", err)
	}

	userDataRes := &requests.UserDataRes{
		ID:              user.ID,
		Nickname:        user.Nickname,
		CareerDirection: user.UserDetail.CareerDirection,
		HomePage:        user.UserDetail.HomePage,
		Signature:       user.UserDetail.Signature,
	}

	for _, tag := range user.Tags {
		userDataRes.UserTags = append(userDataRes.UserTags, tag.Name)
	}

	// 查询所有 Tags 的名字
	var allTags []models.Tag

	if err := db.Select("name").Find(&allTags).Error; err != nil {
		return nil, fmt.Errorf("SelectPersonData -> %s", err)
	}

	for _, tag := range allTags {
		userDataRes.AllTagNames = append(userDataRes.AllTagNames, tag.Name)
	}

	// 获取用户头像图片
	/*// 使用 strconv.ParseUint 将字符串解析为 uint64 类型
	value, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("SelectPersonData -> %s", err)
	}

	// 将 uint64 类型转换为 uint
	uintValue := uint(value)*/
	images, err := controllers.GetImagesControllers("用户", userID)
	if err == nil {
		// 获取图片路径
		for _, image := range *images {
			userDataRes.Path = image.Path
		}
	} else {
		userDataRes.Path = internal_utils.UserDefaultImage
	}

	return userDataRes, nil

}

// UserAccountResponse 响应用户账号设置
func UserAccountResponse(userID uint, db *gorm.DB) (*requests.UserAccountRes, error) {
	var user models.User

	err := db.Preload("UserDetail").First(&user, userID).Error
	if err != nil {
		return nil, fmt.Errorf("UserAccountResponse -> %s", err)
	}

	userAccountRes := &requests.UserAccountRes{
		ID:         user.ID,
		Email:      user.Email,
		BlogLink:   user.UserDetail.BlogLink,
		WeiboLink:  user.UserDetail.WeiboLink,
		GithubLink: user.UserDetail.GithubLink,
		Password:   user.Password,
	}

	return userAccountRes, nil

}

// UserPrivateSetResponse 响应用户私信设置
func UserPrivateSetResponse(userID uint, db *gorm.DB) (*requests.UserPrivateSettingsRes, error) {
	var user models.User

	err := db.Select("PrivateSettings").First(&user, userID).Error
	if err != nil {
		return nil, fmt.Errorf("UserPrivateSetResponse -> %s", err)
	}

	userPrivateSetRes := &requests.UserPrivateSettingsRes{
		ID:              userID,
		PrivateSettings: user.PrivateSettings,
	}

	return userPrivateSetRes, nil

}
