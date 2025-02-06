package controllers

import (
	"forum/internal/message/logics"
	"github.com/gin-gonic/gin"
)

// MessagePushCtrl 向用户实时发送更新数据
func MessagePushCtrl(c *gin.Context) {

	// 设置sse响应的响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream") // 标志了该响应为事件流类型
	c.Writer.Header().Set("Cache-Control", "no-cache")         // 提示用户不要缓存响应
	c.Writer.Header().Set("Connection", "keep-alive")          // 保持连接不断开，以便持续发送事件
	//c.Writer.Header().Set("Access-Control-Allow-Origin", "*")  // 设置跨域资源共享头，允许所有域访问该资源
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://192.168.10.7:9901") // 设置跨域资源共享头，允许所有域访问该资源

	// 创建用户的消息通道 监听新消息并发送到客户端
	logics.NewMessageChan(c)

	// 监听通道中的新评论消息
	// 在外面实现
}
