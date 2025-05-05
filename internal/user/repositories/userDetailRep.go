package repositories

import (
	"errors"
	"fmt"
	"forum/internal/internalPkg/internalUtils"
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
		// 如果没找到，就可以更改
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			return fmt.Errorf("QueryPersonEmail -> 查询 Email异常 -> %s", err)
		}
	}
	// 如果查出来的是自己旧的Email，那也可以更改，否则就不能更改
	if user.ID == userAccountReq.ID {
		return nil
	}

	return fmt.Errorf("该 Email 已经被其他人使用")
}

// UserDataRequest 更新用户个人资料
func UserDataRequest(userId uint, userDataReq *requests.UserDataReq, db *gorm.DB) error {

	var user models.User

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("UserDataRequest -> 开启事务失败 -> %s", tx.Error)
	}

	// 查询该用户是否存在
	err := tx.Model(&models.User{}).Where("id = ?", userId).First(&user).Error
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
	err = tx.Where("user_id", userId).First(&userDetail).Error

	if err != nil {
		userDetail.ID = userId
		userDetail.CareerDirection = userDataReq.CareerDirection
		userDetail.HomePage = userDataReq.HomePage
		userDetail.Signature = userDataReq.Signature
		// return fmt.Errorf("UserDataRequest -> 用户详情表中用户不存在 -> %s", err)
		err = tx.Create(&userDetail).Error
		if err != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("UserDataRequest -> 用户详情表中数据插入失败 -> %s", err)
		}
	} else {
		// 使用 Map 更新特定字段，如果 UserDataReq 结构体字段与数据库字段不一致时
		updates := map[string]interface{}{
			"career_direction": userDataReq.CareerDirection,
			"home_page":        userDataReq.HomePage, // 注意这里要使用数据库中的列名
			"signature":        userDataReq.Signature,
		}

		// 更新用户详情表中相应的字段
		err = tx.Model(&userDetail).Select("CareerDirection", "HomePage", "Signature").Updates(updates).Error
		if err != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("UserDataRequest -> 用户详情表中没有更新任何记录 -> %s", err)
		}
	}
	//// 更新用户标签
	//var existingTags []models.Tag
	//err = tx.Where("name IN ?", userDataReq.UserTags).Find(&existingTags).Error
	//if err != nil {
	//	tx.Rollback() // 回滚事务
	//	return fmt.Errorf("UserDataRequest -> 更新用户标签失败 -> %s", err)
	//}
	//
	//// 找到所有传递过来的标签的ID
	//var tagIDs []uint
	//for _, tag := range existingTags {
	//	tagIDs = append(tagIDs, tag.ID)
	//}

	// 更新用户标签

	// 查询该用户旧的标签
	var ut []models.UserTag
	err = tx.Model(&models.UserTag{}).Where("user_id = ?", userId).Find(&ut).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("UserDataRequest -> 用户详情表中没有更新任何记录 -> %s", err)
	}
	fmt.Println(len(ut))

	// 删除用户旧的标签
	for _, userTag := range ut {
		//fmt.Println("-----------------------------------------------删除")
		result := tx.Delete(&userTag)
		if result.Error != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("UserDataRequest -> 清除旧的用户标签异常 -> %s", err)
		} else if result.RowsAffected == 0 {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("UserDataRequest -> 清除旧的用户标签失败 -> %s", err)
		}
	}
	//fmt.Println("*************************************************1111111")

	// 查询用户选择的标签的id
	var tagIDs []uint
	err = tx.Model(models.Tag{}).Where("name IN ?", userDataReq.UserTags).Pluck("id", &tagIDs).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("UserDataRequest -> 更新用户标签失败 -> %s", err)
	}

	// 添加新的用户标签关联
	for _, tagId := range tagIDs {
		//fmt.Println("-----------------------------------------------------------插入")
		userTag := models.UserTag{
			UserID: userId,
			TagID:  tagId,
		}

		result := tx.Create(&userTag)
		//fmt.Println("---------->", userTag)
		if result.Error != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("UserDataRequest -> 添加新的用户标签关联失败 -> %s", err)
		} else if result.RowsAffected == 0 {
			fmt.Println("------------>无数据插入")
		}
	}

	//// 添加新的用户标签关联
	//for _, tagID := range tagIDs {
	//	err := tx.Create(&models.UserTag{UserID: userId, TagID: tagID}).Error
	//	if err != nil {
	//		tx.Rollback() // 回滚事务
	//		return fmt.Errorf("UserDataRequest -> 添加新的用户标签关联失败 -> %s", err)
	//	}
	//}

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
	if userAccountReq.Password != "" {
		// 对密码进行加密
		hashPassword, err := internalUtils.HashPassword(userAccountReq.Password)
		if err != nil {
			return fmt.Errorf("UserAccountRequest -> 密码加密失败 -> %s", err)
		}

		// 更新 User 表中的, password
		err = tx.Model(&user).Updates(map[string]interface{}{
			"password": hashPassword,
		}).Error
		if err != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("UserAccountRequest -> 更新 User 表中的 email , password -> %s", err)
		}
	}
	/*else {
		// 更新 User 表中的 email , password
		err = tx.Model(&user).Updates(map[string]interface{}{
			"email": userAccountReq.Email,
		}).Error
		if err != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("UserAccountRequest -> 更新 User 表中的 email -> %s", err)
		}
	}*/

	// 查询该用户的外键是否存在
	var userDetail models.UserDetail
	err = tx.Where("user_id = ?", userAccountReq.ID).First(&userDetail).Error

	if err != nil {
		// return fmt.Errorf("UserAccountRequest -> 用户详情表中用户不存在 -> %s", err)
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
func UserPrivateSetRequest(userID uint, userPrivateSetReq *requests.UserPrivateSettingsReq, db *gorm.DB) error {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("UserPrivateSetRequest -> 开启事务失败 -> %s", tx.Error)
	}

	var user models.User
	fmt.Println("**************----->", userID)
	// 查询该用户是否存在
	err := tx.Model(&models.User{}).Where("id = ?", userID).First(&user).Error
	fmt.Println("^^^^^^^^^^^^^-->", err)
	if err != nil {
		// 数据库表中还没有该用户的数据，直接插入即可
		user.ID = userID
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

	err := db.Preload("UserDetail").First(&user, userID).Error
	if err != nil {
		return nil, fmt.Errorf("UserDataResponse -> 获取用户信息异常 -> %s", err)
	}

	userDataRes := &requests.UserDataRes{
		ID:              user.ID,
		Nickname:        user.Nickname,
		CareerDirection: user.UserDetail.CareerDirection,
		HomePage:        user.UserDetail.HomePage,
		Signature:       user.UserDetail.Signature,
	}

	// 查询用户选择的标签id
	var newTag []uint
	err = db.Model(&models.UserTag{}).Where("user_id = ?", userID).Pluck("tag_id", &newTag).Error
	if err != nil {
		return nil, fmt.Errorf("UserDataResponse -> 查询用户选择的标签异常 -> %s", err)
	}

	// 查询用户选择的标签的名字
	var tagName []string
	err = db.Model(&models.Tag{}).Where("id IN ?", newTag).Pluck("name", &tagName).Error
	if err != nil {
		return nil, fmt.Errorf("UserDataResponse -> 查询用户选择的标签的名字异常 -> %s", err)
	}

	userDataRes.UserTags = tagName

	//for _, tag := range user.Tags {
	//	userDataRes.UserTags = append(userDataRes.UserTags, tag.Name)
	//	fmt.Println("------------------>", tag.ID, tag.Name)
	//}

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
	images, err := internalUtils.GetImages(db, globals.UserHome, userID)
	if err == nil {
		// 获取图片路径
		for _, path := range *images {
			userDataRes.Path = path
		}
	} else {
		userDataRes.Path = internalUtils.UserDefaultImage
	}

	return userDataRes, nil

}

// UserAccountResponse 响应用户账号设置
func UserAccountResponse(userID string, db *gorm.DB) (*requests.UserAccountRes, error) {
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
		Password:   "",
		Nickname:   user.Nickname,
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
