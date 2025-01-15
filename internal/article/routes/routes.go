package routes

import (
	"forum/internal/article/controllers"
	"forum/pkg/token"
	"github.com/gin-gonic/gin"
)

// Article
// @Description: 文章
// @param        e *gin.Engine
func Article(e *gin.Engine) {

	// 搜索
	articleGroup := e.Group("/article")
	// token 校验
	articleGroup.Use(token.AuthMiddleware())
	{
		// 文章搜索框
		articleGroup.GET("/search_box", controllers.ArticleSearchCtrl)
		// 编辑界面
		articleGroup.GET("/edit", controllers.ArticleEditCtrl)
		// 发布文章
		articleGroup.POST("/publish", controllers.ArticlePublishCtrl)
		// 获取文章列表
		articleGroup.POST("/get_list", controllers.ArticleListCtrl)
		// 封禁文章
		articleGroup.GET("/ban", controllers.ArticleBanCtrl)
		// 删除文章
		articleGroup.DELETE("/delete", controllers.DeleteArticlesCtrl)
		// 获取文章详情
		articleGroup.GET("/detail", controllers.ArticleDetailCtrl)
		// 点赞
		articleGroup.POST("/like", controllers.LikeArticleCtrl)
		// 收藏
		articleGroup.POST("/collection", controllers.CollectionCtrl)

		// 获取标签下的文章
		articleGroup.GET("/get_article_by_tag", controllers.GetArticlesByTagCtrl)
	}

	// token 校验
	e.Use(token.AuthMiddleware())
	// 会员中心 获取用户文章或收藏列表
	e.GET("/get_type_data", controllers.GetUserArticleOrCollectionCtrl)

}

// Workplace
// @Description: 工作台路由
// @param        e *gin.Engine
// @Author tianjiajie 2025-01-15 09:11:06
func Workplace(e *gin.Engine) {
	workplaceGroup := e.Group("/workplace")
	// token 校验
	workplaceGroup.Use(token.AuthMiddleware())
	{
		// 查询近两周文章发布数量
		workplaceGroup.GET("/article_sum", controllers.GetTwoWeeksArticleSumCtrl)
		// 查询前五篇热门文章数据
		workplaceGroup.GET("/hot_articles", controllers.GetHotArticleCtrl)
		// 前五个热门标签的文章量
		workplaceGroup.GET("/pielist", controllers.GetHotTagsCtrl)

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
	r.POST("/delete", controllers.DeleteCommentCtrl)

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
