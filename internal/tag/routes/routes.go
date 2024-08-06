package routes

import (
	"forum/internal/tag/controllers"
	"github.com/gin-gonic/gin"
)

func RouterInit(e *gin.Engine) {
	// 分组
	r := e.Group("/tag")

	// 标签人数更新
	r.POST("/fan_count", controllers.UpdateTagUserCount)
}
