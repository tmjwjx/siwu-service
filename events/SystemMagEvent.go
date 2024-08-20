package events

type SystemMsgEvent struct {
	Name string
}

func NewSystemMsgEvent(name string) *SystemMsgEvent {
	s := &SystemMsgEvent{
		Name: name,
	}
	return s
}

func (s *SystemMsgEvent) GetEventName() string {
	return s.Name
}