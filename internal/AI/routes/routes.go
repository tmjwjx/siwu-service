package routes

import (
	"forum/internal/AI/controllers"
	"github.com/gin-gonic/gin"
)

func AIRouters(e *gin.Engine) {
	e.POST("/AI/codeExplain", controllers.GetCodeExplain)

	articleRouter := e.Group("/AI/extract")
	{
		// 第一次获取文章的摘要、总结、标签
		articleRouter.POST("/get_info_first", controllers.GetArticleInfoFirstCtrl)

		// 将文章的ID保存到相应的记录中
		articleRouter.POST("/save_article_id", controllers.SaveArticleIDCtrl)

		// 非首次获取文章的摘要、总结、标签
		articleRouter.GET("/get_info", controllers.GetArticleInfoCtrl)

		// 删除文章相关信息
		articleRouter.DELETE("/del_info", controllers.DelArticleInfoCtrl)
	}
}
