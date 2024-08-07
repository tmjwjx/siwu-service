package repositories

import (
	"fmt"
	"forum/internal/image/controllers"
	"forum/internal/models"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

// QueryPersonEmail 查询Email
func QueryPersonEmail(u *requests.UserResponse, c *gin.Context) *gorm.DB {

	// 查询 Email 是否唯一
	var user models.User
	result := globals.DB.Where("email = ?", u.User.Email).First(&user)
	return result
}

// UpdatePersonData 更新用户信息
func UpdatePersonData(u *requests.UserResponse, c *gin.Context) *gorm.DB {

	// 将前端传过来的 user 文本类数据插入到数据库中
	result2 := globals.DB.Create(&u)
	return result2
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
	// 使用 strconv.ParseUint 将字符串解析为 uint64 类型
	value, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("SelectPersonData -> %s", err)
	}

	// 将 uint64 类型转换为 uint
	uintValue := uint(value)
	images, err := controllers.GetImagesControllers("用户", uintValue)
	if err != nil {
		return nil, fmt.Errorf("SelectPersonData -> %s", err)
	}
	// 获取图片路径
	var path string
	for _, image := range *images {
		path = image.Path
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
