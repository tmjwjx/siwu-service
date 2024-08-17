package routes

import (
	"forum/internal/article/controllers"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
)

func Search() {
	// 搜索
	userGroup := globals.Router.Group("/search")
	{
		userGroup.GET("/query", controllers.ArticleSearchCtrl)
	}

}

func Publish() {

	// 搜索
	articleGroup := globals.Router.Group("/article")
	{
		// 编辑界面
		articleGroup.GET("/edit", controllers.ArticleEditCtrl)
		// 发布文章
		articleGroup.POST("/publish", controllers.ArticlePublishCtrl)
	}

}

// Comment 评论
func Comment(e *gin.Engine) {

	r := e.Group("/comment")

	// 保存评论
	r.POST("/create", controllers.InsertCommentCtrl)

	// 返回顶级评论
	r.GET("/top_level", controllers.GetTopLevelCommentsCtrl)

	// 返回评论回复
	r.GET("/replies", controllers.GetRepliesRep2Ctrl)

	// 删除评论
	r.DELETE("/delete", controllers.DeleteCommentCtrl)

	// 更新点赞的数量
	r.POST("/praise", controllers.UpdatePraiseCountCtrl)
}
