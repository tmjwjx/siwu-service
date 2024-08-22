package routes

import (
	"forum/internal/article/controllers"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
)

// Search 搜索文章
func Search() {
	// 搜索
	userGroup := globals.Router.Group("/search")
	{
		userGroup.GET("/query", controllers.ArticleSearchCtrl)
	}

}

// Publish 发布文章
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

	// 前台
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

	// 后台
	r2 := e.Group("/backstage_comment")

	// 展示评论列表
	r2.GET("/list", controllers.ShowCommentsListCtrl)

	// 添加评论
	r2.POST("/add", controllers.AddCommentCtrl)

	// 删除评论
	r2.DELETE("/delete", controllers.BsDeleteCommentCtrl)

	// 批量删除评论
	r2.DELETE("/batch_delete", controllers.BatchDelTagCtrl)

	// 更新评论
	r2.POST("/update", controllers.UpdateCommentCtrl)

	// 查询某个用户的全部评论
	r2.GET("/query", controllers.QueryCommentCtrl)

}
