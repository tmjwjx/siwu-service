package controllers

import (
	"forum/internal/user/logics"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
)

// 注册、登陆、验证码

// Register 注册
func Register(c *gin.Context) {
	userLogic := logics.NewUserLogic(globals.DB, c)

	// 业务逻辑
	appErr := userLogic.Register()
	if appErr != nil {
		response.Failed(c, appErr)
		return
	}

	// 调用默认的成功
	response.Success(c, response.StatusOkData)
}
