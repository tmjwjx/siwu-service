package routes

import (
	"forum/internal/tag/controllers"
	"github.com/gin-gonic/gin"
)

func Tag(e *gin.Engine) {
	// 分组
	r := e.Group("/tag")

	// 更新标签关注人数
	r.POST("/fan_count", controllers.UpdateTagUserCount)

	// 更新标签文章数量
	r.POST("/article_count", controllers.UpdateTagArticleCount)

	// 更新标签热度
	r.POST("/heat", controllers.UpdateTagHeat)
}
