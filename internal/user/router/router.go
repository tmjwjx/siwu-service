package router

import (
	"forum/internal/user/controllers"
	"github.com/gin-gonic/gin"
)

// RouterInit 用户路由
func RouterInit(e *gin.Engine) {
	// 分组
	r := e.Group("/user")

	r.POST("/register", controllers.Register)

}
