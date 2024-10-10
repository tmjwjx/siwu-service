package observer

import (
	"fmt"
	"forum/pkg/globals"
)

type SystemMsgObserver struct {
	userId string
}

// NewSystemMsgObserver
// @Description: 创建一个新的系统消息观察者
// @param        userId string
// @return       *SystemMsgObserver
// @Author tianjiajie 2024-10-05 20:33:20
func NewSystemMsgObserver(userId string) *SystemMsgObserver {
	return &SystemMsgObserver{
		userId: userId,
	}
}

// Update
// @Description: 执行观察者的更新方法
// @receiver     u
// @param        event event.Subject
// @Author tianjiajie 2024-10-09 16:03:36
func (u *SystemMsgObserver) Update(data string) {
	//if u.userId != userId {
	//	return
	//}

	// 消息通知事件 推送消息给指定用户
	notifyChan, exists := globals.SubscriberChannels[u.userId]
	fmt.Println("exist", exists)
	if exists {
		notifyChan <- data // 推送消息给指定用户的通道
		globals.Log.Infof("用户 %s 收到消息 %s", u.userId, data)
	}
}

//// PushMessage
//// @Description: 推送消息
//// @param        userId string
//// @param        eventName string
//// @Author tianjiajie 2024-10-05 20:35:27
//func PushMessage(userId, eventName string) {
//	systemMsgObserver := NewSystemMsgObserver(userId)
//	fmt.Println(systemMsgObserver)
//	systemMsgEvent := events.NewSystemMsgEvent(eventName)
//	fmt.Println(systemMsgEvent)
//
//	globals.BusDispatcher.Register(systemMsgEvent, systemMsgObserver)
//	fmt.Println("testMap", globals.BusDispatcher.Event[eventName][0])
//
//	err := globals.BusDispatcher.Dispatch(systemMsgEvent)
//	if err != nil {
//		globals.Log.Error(err)
//	}
//	fmt.Println("hello world    1")
//}