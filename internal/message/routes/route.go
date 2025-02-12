package routes

import (
	"forum/internal/message/controllers"
	"forum/pkg/token"
	"github.com/gin-gonic/gin"
)

func Message(e *gin.Engine) {
	
	messageGroup := e.Group("/message")
	{
		messageGroup.Use(token.AuthMiddleware())
		// 推送消息
		messageGroup.GET("/sse", controllers.MessagePushCtrl)
		// 点赞消息
		messageGroup.GET("/like", controllers.LikeMessageCtrl)
		// 收藏消息
		messageGroup.GET("/collection", controllers.CollectionMessageCtrl)
		// 关注消息
		messageGroup.GET("/follow", controllers.FollowMessageCtrl)
		// 评论消息
		messageGroup.GET("/comment", controllers.CommentMesCtrl)
		// 评论点赞消息
		messageGroup.GET("/comment_like", controllers.CommentLikeMesCtrl)
		
		// 未读点赞消息数量
		messageGroup.GET("/like_unread", controllers.LikeUnreadCtrl)
		// 未读收藏消息数量
		messageGroup.GET("/like_unread", controllers.LikeUnreadCtrl)
		// 未读评论消息数量
		messageGroup.GET("/like_unread", controllers.LikeUnreadCtrl)
		// 未读关注消息数量
		messageGroup.GET("/like_unread", controllers.LikeUnreadCtrl)
		// 未读评论点赞消息数量
		messageGroup.GET("/like_unread", controllers.LikeUnreadCtrl)
	}
	
}
