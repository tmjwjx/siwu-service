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

func NewSystemMsgObserver(userId string) *SystemMsgObserver {
	return &SystemMsgObserver{
		userId: userId,
	}
}

func (u *SystemMsgObserver) Process(event event.Event) {
	fmt.Println("Process 启动")
	fmt.Println(u.userId)
	switch ev := event.(type) {
	case *events.SystemMsgEvent:
		notifyChan, exists := globals.SubscriberChannels[u.userId]
		fmt.Println("exist", exists)
		if exists {
			fmt.Println("notifyChan存在")
			notifyChan <- ev.Name // 推送消息给指定用户的通道
		}
	}
}

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