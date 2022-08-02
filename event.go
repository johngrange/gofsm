package fsm

type eventImpl struct {
	name string
	data interface{}
}

func NewEvent(name string, data interface{}) Event {
	return &eventImpl{
		name: name,
		data: data,
	}
}

func (ev *eventImpl) Name() string {
	return ev.name
}

func (ev *eventImpl) Data() interface{} {
	return ev.data
}
