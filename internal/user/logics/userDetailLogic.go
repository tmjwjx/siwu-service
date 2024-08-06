package logics

import (
	"forum/internal/image/controllers"
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"github.com/gin-gonic/gin"
)

// PersonalDataLogic 将用户信息存入数据库
func PersonalDataLogic(u *requests.User, c *gin.Context) (error, int) {

	// 验证Email是否唯一
	result := repositories.QueryPersonEmail(u, c)
	if result.Error == nil {
		// 查询到数据库已经存在了该 Email，返回错误
		return result.Error, 400
	}

	// 运行到这一步，已经证明数据库中还没有该 Email ，接下来将其信息存入数据库中
	result2 := repositories.UpdatePersonData(u, c)
	if result2.Error != nil {
		// 插入数据失败，返回错误
		return result2.Error, 500
	}

	// 将图片文件的路径相关信息存入数据库
	err, status := controllers.UploadHandlerControllers(c, 1, u.ID)
	if err != nil {
		return err, status
	}
	return nil, 200
}

// ResponsePersonDateLogic 将用户信息响应给前端
func ResponsePersonDateLogic(userID string, c *gin.Context) (*requests.User, error) {
	user, err := repositories.SelectPersonData(userID)
	return user, err
}
