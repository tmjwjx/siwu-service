package routes

import (
	"forum/internal/user/controllers"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
)

// RouterInit 用户路由
func RouterInit(e *gin.Engine) {
	// 分组
	r := e.Group("/user")

	r.POST("/register", controllers.Register)

	// 上传个人资料
	globals.Router.POST("/form_personal_data", controllers.PersonalDataHandler)

	// 前端获取用户资料
	globals.Router.GET("/form_personal_data/:id", controllers.ResponsePersonDate)
}
