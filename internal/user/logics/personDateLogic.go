package logics

import (
	"forum/internal/image/controllers"
	"forum/internal/user/repositorys"
	"forum/internal/user/requests"
	"github.com/gin-gonic/gin"
)

func PersonalDataLogic(u *requests.User, c *gin.Context) {

	// 验证Email是否唯一
	repositorys.Query(u, c)

	// 运行到这一步，已经证明数据库中还没有该 Email ，接下来将其存入数据库中
	repositorys.Update(u, c)

	// 将图片文件的路径存入数据库
	controllers.UploadHandlerControllers(c)
}
