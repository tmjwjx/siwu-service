package routes

import (
	"forum/internal/AI/controllers"
	"github.com/gin-gonic/gin"
)

func AIRouters(e *gin.Engine) {
	articleRouter := e.Group("/AI/extract")
	{
		articleRouter.GET("/get_info_first", controllers.GetArticleInfoFirstCtrl)
		articleRouter.POST("/save_article_id", controllers.SaveArticleIDCtrl)
		articleRouter.GET("/get_info", controllers.GetArticleInfoCtrl)
		articleRouter.DELETE("/del_info", controllers.DelArticleInfoCtrl)
	}
}
