package repositories

import (
	"forum/internal/models"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Query(u *requests.User, c *gin.Context) {

	// 查询 Email 是否唯一
	var user models.User
	result := globals.DB.Where("email = ?", u.Email).First(&user)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "邮件已经存在"})
		return
	}
}

func Update(u *requests.User, c *gin.Context) {

	// 将前端传过来的 user 文本类数据插入到数据库中
	result2 := globals.DB.Create(&u)
	if result2.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "用户数据保存失败"})
	}
	c.JSON(http.StatusOK, gin.H{"msg": "用户数据保存成功"})
}

func SelectPersonData(userID string) (*requests.User, error) {
	var user requests.User

	// 预加载 UserDetail
	if err := globals.DB.Preload("UserDetail").First(&user, userID).Error; err != nil {
		return &user, err
	}
	return &user, nil
}
