package event

// Observer 观察者
type Observer interface {
	Update(data string) // 执行业务
}

//// Subject 事件
//type Subject interface {
//	//GetEventName() string // 获取事务名称
//}

//// Dispatcher 调度员
//type Dispatcher struct {
//	Event map[string][]Observer
//}

//// NewDispatcher 创建调度员
//func NewDispatcher() *Dispatcher {
//	return &Dispatcher{
//		Event: make(map[string][]Observer),
//	}
//}

//// Register 调度连接，连接事件名称和观察者
//func (e *Dispatcher) Register(event Subject, observer Observer) {
//	fmt.Println(event.GetEventName())
//	fmt.Println(observer)
//	e.Event[event.GetEventName()] = append(e.Event[event.GetEventName()], observer)
//}

//// Dispatch  事件发生 通知事件的观察者
//func (e *Dispatcher) Dispatch(event Subject) error {
//	observer, ok := e.Event[event.GetEventName()]
//	if ok {
//		for i, v := range observer {
//			fmt.Println("i = ", i)
//			fmt.Println(event)
//			v.Update(event)
//		}
//		fmt.Println("hello world    2")
//		return nil
//	}
//	return errors.New("未知错误")
//}