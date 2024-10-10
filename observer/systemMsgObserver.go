package observer

import (
	"fmt"
	"forum/events"
	"forum/pkg/event"
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

// Process
// @Description: 处理事件
// @receiver     u
// @param        event event.Event
// @Author tianjiajie 2024-10-05 20:33:26
func (u *SystemMsgObserver) Process(event event.Event) {
	//fmt.Println("Process 启动")
	//fmt.Println(u.userId)
	switch ev := event.(type) {
	// 消息通知事件 推送消息给指定用户
	case *events.SystemMsgEvent:
		notifyChan, exists := globals.SubscriberChannels[u.userId]
		fmt.Println("exist", exists)
		if exists {
			fmt.Println("notifyChan存在")
			notifyChan <- ev.Name // 推送消息给指定用户的通道
		}
	}
}

// PushMessage
// @Description: 推送消息
// @param        userId string
// @param        eventName string
// @Author tianjiajie 2024-10-05 20:35:27
func PushMessage(userId, eventName string) {
	systemMsgObserver := NewSystemMsgObserver(userId)
	fmt.Println(systemMsgObserver)
	systemMsgEvent := events.NewSystemMsgEvent(eventName)
	fmt.Println(systemMsgEvent)

	globals.BusDispatcher.Register(systemMsgEvent, systemMsgObserver)
	fmt.Println("testMap", globals.BusDispatcher.Event[eventName][0])

	err := globals.BusDispatcher.Dispatch(systemMsgEvent)
	if err != nil {
		globals.Log.Error(err)
	}
	fmt.Println("hello world    1")
}