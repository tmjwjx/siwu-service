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
		messageGroup.GET("/push", controllers.MessagePushCtrl)
		// 点赞消息
		messageGroup.GET("/like", controllers.LikeMessageCtrl)
		// 收藏消息
		messageGroup.GET("/collection", controllers.CollectionMessageCtrl)
	}

}