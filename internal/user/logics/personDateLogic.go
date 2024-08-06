package logics

import (
	"forum/internal/image/controllers"
	"forum/internal/user/repositorys"
	"forum/internal/user/requests"
	"github.com/gin-gonic/gin"
)

// 将用户信息存入数据库
func PersonalDataLogic(u *requests.User, c *gin.Context) {

	// 验证Email是否唯一
	repositorys.Query(u, c)

	// 运行到这一步，已经证明数据库中还没有该 Email ，接下来将其信息存入数据库中
	repositorys.Update(u, c)

	// 将图片文件的路径相关信息存入数据库
	controllers.UploadHandlerControllers(c, 1, u.ID)
}

// 将用户信息响应给前端
func ResponsePersonDateLogic(userID string, c *gin.Context) (*requests.User, error) {
	user, err := repositorys.SelectPersonData(userID)
	return user, err
}
