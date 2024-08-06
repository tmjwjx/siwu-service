package repositories

import (
	"forum/internal/models"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// QueryPersonEmail 查询Email
func QueryPersonEmail(u *requests.User, c *gin.Context) *gorm.DB {

	// 查询 Email 是否唯一
	var user models.User
	result := globals.DB.Where("email = ?", u.Email).First(&user)
	return result
}

// UpdatePersonData 更新用户信息
func UpdatePersonData(u *requests.User, c *gin.Context) *gorm.DB {

	// 将前端传过来的 user 文本类数据插入到数据库中
	result2 := globals.DB.Create(&u)
	return result2
}

// SelectPersonData 查询用户信息
func SelectPersonData(userID string) (*requests.User, error) {
	var user requests.User

	// 预加载 UserDetail
	if err := globals.DB.Preload("UserDetail").First(&user, userID).Error; err != nil {
		return &user, err
	}
	return &user, nil
}
