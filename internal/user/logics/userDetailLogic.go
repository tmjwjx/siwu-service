package logics

import (
	"forum/internal/image/controllers"
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"github.com/gin-gonic/gin"
)

// PersonalDataLogic 将用户信息存入数据库
func PersonalDataLogic(u *requests.User, c *gin.Context) {

	// 验证Email是否唯一
	repositories.Query(u, c)

	// 运行到这一步，已经证明数据库中还没有该 Email ，接下来将其信息存入数据库中
	repositories.Update(u, c)

	// 将图片文件的路径相关信息存入数据库
	controllers.UploadHandlerControllers(c, 1, u.ID)
}

// 将用户信息响应给前端
func ResponsePersonDateLogic(userID string, c *gin.Context) (*requests.User, error) {
	user, err := repositories.SelectPersonData(userID)
	return user, err
}
