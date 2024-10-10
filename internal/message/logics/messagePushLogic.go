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
	id, _ := c.Get("id")
	userId := id.(string)
	notifyChan := make(chan string)
	globals.SubscriberChannels[userId] = notifyChan
	//notifyChan := globals.SubscriberChannels[userId]

	//// 创建一个观察者
	//observer.NewSystemMsgObserver(userId)
	//// 注册观察者到事件
	//globals.SystemMsgSubject.RegisterObserver("systemMsg", observer.NewSystemMsgObserver(userId))

	for {
		select {
		case message := <-notifyChan:
			_, err := fmt.Fprintf(c.Writer, "data: %s\n\n", message)
			if err != nil {
				continue
			}
			c.Writer.Flush()
		case <-c.Done():
			// 关闭通道
			close(notifyChan)

			// 删除用户的通道
			delete(globals.SubscriberChannels, userId)

			//// 取消注册观察者
			//globals.SystemMsgSubject.UnRegisterObserver("systemMsg", observer.NewSystemMsgObserver(userId))
			return
		}
	}
}