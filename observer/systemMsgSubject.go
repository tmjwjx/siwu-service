package observer

import (
	"forum/pkg/event"
	_ "forum/pkg/event"
	"github.com/asaskevich/EventBus"
)

type SystemMsgSubject struct {
	EventBus EventBus.Bus
	//name string
}

// NewSystemMsgSubject
// @Description: 创建一个新的系统消息主题
// @return       *SystemMsgSubject
// @Author tianjiajie 2024-10-09 15:55:56
func NewSystemMsgSubject() *SystemMsgSubject {
	return &SystemMsgSubject{
		EventBus: EventBus.New(),
	}
}

// RegisterObserver
// @Description: 注册观察者到特定事件
// @receiver     s
// @param        event string
// @param        observer event.Observer
// @Author tianjiajie 2024-10-09 15:59:49
func (s *SystemMsgSubject) RegisterObserver(event string, observer event.Observer) {
	_ = s.EventBus.Subscribe(event, observer.Update)
}

// UnRegisterObserver
// @Description: 取消注册观察者
// @receiver     s
// @param        event string
// @param        observer event.Observer
// @Author tianjiajie 2024-10-09 16:01:34
func (s *SystemMsgSubject) UnRegisterObserver(event string, observer event.Observer) {
	_ = s.EventBus.Unsubscribe(event, observer.Update)
}

// Notify
// @Description: 发布事件给所有观察者
// @receiver     s
// @param        event string
// @param        data string
// @Author tianjiajie 2024-10-09 16:01:46
func (s *SystemMsgSubject) Notify(event string, data string, userId string) {
	// 发布事件 event 给所有观察者 data为发送内容
	s.EventBus.Publish(event, data)
}

//func (s *SystemMsgSubject) GetEventName() string {
//	return s.Name
//}