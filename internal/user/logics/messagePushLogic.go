package logics

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
)

// TriggerEvent 模拟外部触发消息推送的函数
func TriggerEvent(userId string, message string) {

	notifyChan, exists := globals.SubscriberChannels[userId]

	if exists {
		notifyChan <- models.Event{Type: message} // 推送消息给指定用户的通道
	}
}

func NewMessageChan(c *gin.Context) {
	// 创建一个通道，用于接收评论的变动消息
	userId := c.Query("user_id")
	notifyChan := globals.SubscriberChannels[userId]

	for {
		select {
		case message := <-notifyChan:
			fmt.Fprintf(c.Writer, "data: %s\n\n", message)
			c.Writer.Flush()
		case <-c.Done():
			close(notifyChan)
			delete(globals.SubscriberChannels, userId)
			return
		}
	}
}