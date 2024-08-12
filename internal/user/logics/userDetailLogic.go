package logics

import (
	"fmt"
	"forum/internal/image/controllers"
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"github.com/gin-gonic/gin"
)

// PersonalDataLogic 将用户信息存入数据库
func PersonalDataLogic(userRequest *requests.UserRequest, c *gin.Context) (error, int) {

	// 验证Email是否唯一
	err := repositories.QueryPersonEmail(userRequest)
	if err == nil {
		// 查询到数据库已经存在了该 Email，返回错误
		return fmt.Errorf("PersonalDataLogic -> %s", "该 Email 已经存在"), 400
	}

	// 运行到这一步，已经证明数据库中还没有该 Email ，接下来将其信息存入数据库中
	err = repositories.UpdatePersonData(userRequest)
	if err != nil {
		return err, 500
	}

	// 将图片文件的路径相关信息存入数据库
	err, status := controllers.UploadImagesControllers(c, "用户", userRequest.User.ID)
	if err != nil {
		return err, status
	}
	return nil, 200
}

// ResponsePersonDateLogic 将用户信息响应给前端
func ResponsePersonDateLogic(userID string) (*requests.UserResponse, error) {
	userResponse, err := repositories.SelectPersonData(userID)
	return userResponse, err
}
