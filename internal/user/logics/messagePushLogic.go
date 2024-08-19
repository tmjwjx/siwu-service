package logics

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
)

type MessageEvent interface {
	PushMessage()
}
type Message struct {
	Name   string
	UserId string
}

func (m Message) PushMessage() {
	notifyChan, exists := globals.SubscriberChannels[m.UserId]
	if exists {
		notifyChan <- models.Event{Type: m.Name} // 推送消息给指定用户的通道
	}
}

// TriggerEvent 模拟外部触发消息推送的函数
func TriggerEvent(message MessageEvent) {
	switch v := message.(type) {
	case Message:
		v.PushMessage()
	default:
		fmt.Printf("Unknown type\n")
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