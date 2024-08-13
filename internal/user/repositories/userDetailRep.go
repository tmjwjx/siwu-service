package repositories

import (
	"fmt"
	"forum/internal/image/controllers"
	"forum/internal/internal_pkg/internal_utils"
	"forum/internal/models"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"strconv"
)

// QueryPersonEmail 查询Email
func QueryPersonEmail(userRequest *requests.UserRequest) error {

	// 查询 Email 是否唯一
	var user models.User
	err := globals.DB.Where("email = ?", userRequest.User.Email).First(&user).Error
	if err != nil {
		return fmt.Errorf("QueryPersonEmail -> %s", err)
	}
	return nil
}

// UpdatePersonData 更新用户信息
func UpdatePersonData(userRequest *requests.UserRequest) error {

	// 将前端传过来的 user 文本类数据插入到数据库中
	// 更新用户信息
	err := globals.DB.Model(&models.User{}).Where("id = ?", userRequest.User.ID).Updates(userRequest.User).Error
	if err != nil {
		return fmt.Errorf("UpdatePersonData -> %s", err)
	}
	// 更新用户详情
	err = globals.DB.Model(&models.UserDetail{}).Where("user_id = ?", userRequest.User.ID).Updates(userRequest.UserDetail).Error
	if err != nil {
		return fmt.Errorf("UpdatePersonData -> %s", err)
	}
	// 更新用户标签
	var existingTags []models.Tag
	err = globals.DB.Where("name IN ?", userRequest.UserTags).Find(&existingTags).Error
	if err != nil {
		return fmt.Errorf("UpdatePersonData -> %s", err)
	}
	// 找到所有传递过来的标签的ID
	var tagIDs []uint
	for _, tag := range existingTags {
		tagIDs = append(tagIDs, tag.ID)
	}
	// 清除旧的用户标签关联
	err = globals.DB.Where("user_id = ?", userRequest.User.ID).Delete(&models.UserTag{}).Error
	if err != nil {
		return fmt.Errorf("UpdatePersonData -> %s", err)
	}
	// 添加新的用户标签关联
	for _, tagID := range tagIDs {
		globals.DB.Create(&models.UserTag{UserID: userRequest.User.ID, TagID: tagID})
	}
	return nil
}

// SelectPersonData 查询用户信息
func SelectPersonData(userID string) (*requests.UserResponse, error) {
	var user models.User

	// 查询 User 和关联的 UserDetail
	if err := globals.DB.Preload("UserDetail").First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("SelectPersonData -> %s", err)
	}
	// // 查询用户关联的 Tags 中的标签名字
	var userTagsNames []string
	if err := globals.DB.Model(&user).Association("Tags").Find(&userTagsNames, "name"); err != nil {
		return nil, fmt.Errorf("SelectPersonData -> %s", err)
	}
	// 查询所有 Tags 的名字
	var allTags []models.Tag
	if err := globals.DB.Select("name").Find(&allTags).Error; err != nil {
		return nil, fmt.Errorf("SelectPersonData -> %s", err)
	}
	var allTagNames []string
	for _, tag := range allTags {
		allTagNames = append(allTagNames, tag.Name)
	}
	// 获取用户头像图片
	// 使用 strconv.ParseUint 将字符串解析为 uint64 类型
	value, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("SelectPersonData -> %s", err)
	}

	// 将 uint64 类型转换为 uint
	uintValue := uint(value)
	images, err := controllers.GetImagesControllers("用户", uintValue)
	var path string
	if err == nil {
		// return nil, fmt.Errorf("SelectPersonData -> %s", err)
		// 获取图片路径
		for _, image := range *images {
			path = image.Path
		}
	} else {
		path = internal_utils.UserDefaultImage
	}

	// 组合数据
	response := &requests.UserResponse{
		User:        user,
		UserDetail:  user.UserDetail,
		UserTags:    userTagsNames,
		AllTagNames: allTagNames,
		Path:        path,
	}
	return response, nil
}
