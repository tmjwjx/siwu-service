package logics

import (
	"fmt"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
)

// NewMessageChan
// @Description: 创建用户的消息通道 监听新消息并发送到客户端
// @param        c *gin.Context
func NewMessageChan(c *gin.Context) {
	// 创建一个通道，用于接收评论的变动消息
	userId := c.Query("user_id")
	notifyChan := make(chan string)
	globals.SubscriberChannels[userId] = notifyChan
	//notifyChan := globals.SubscriberChannels[userId]

	for {
		select {
		case message := <-notifyChan:
			_, err := fmt.Fprintf(c.Writer, "data: %s\n\n", message)
			if err != nil {
				continue
			}
			c.Writer.Flush()
		case <-c.Done():
			close(notifyChan)
			delete(globals.SubscriberChannels, userId)
			return
		}
	}
}