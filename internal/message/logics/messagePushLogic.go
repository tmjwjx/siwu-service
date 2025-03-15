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
	// 获取用户id
	id, _ := c.Get("id")
	globals.Log.Infof("用户 %s 连接成功", id)
	userId := fmt.Sprintf("%d", id)
	
	// 将用户的通道存入全局变量
	notifyChan := make(chan string, 10)
	globals.SubscriberChannels[userId] = notifyChan
	//notifyChan := globals.SubscriberChannels[userId]
	
	//// 创建一个观察者
	//observer.NewSystemMsgObserver(userId)
	//// 注册观察者到事件
	//globals.SystemMsgSubject.RegisterObserver("systemMsg", observer.NewSystemMsgObserver(userId))
	
	for {
		select {
		case message := <-notifyChan:
			// "data: %s\n\n" 是SSE协议发送的固定格式
			//data := gin.H{"type": message}
			//response, err := json.Marshal(data)
			//if err != nil {
			//	continue
			//}
			_, err := fmt.Fprintf(c.Writer, "data: %s\n\n", message)
			//_, err = fmt.Fprintf(c.Writer, "data: %s\n\n", message)
			if err != nil {
				globals.Log.Errorf("发送消息%s失败：%s", message, err)
				continue
			}
			c.Writer.Flush()
		case <-c.Done():
			// 确保通道的安全性
			if _, ok := globals.SubscriberChannels[userId]; ok {
				// 关闭通道
				close(notifyChan)
				
				// 删除用户的通道
				delete(globals.SubscriberChannels, userId)
			}
			
			//// 取消注册观察者
			//globals.SystemMsgSubject.UnRegisterObserver("systemMsg", observer.NewSystemMsgObserver(userId))
			return
		}
	}
}
