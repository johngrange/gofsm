package fsm

type eventImpl struct {
	name   string
	data   interface{}
	labels []string
}

func NewEvent(name string, data interface{}, labels ...string) Event {
	ev := &eventImpl{
		name:   name,
		data:   data,
		labels: []string{},
	}
	ev.labels = append(ev.labels, labels...)
	return ev
}

func (ev *eventImpl) Name() string {
	return ev.name
}

func (ev *eventImpl) Data() interface{} {
	return ev.data
}

func (ev *eventImpl) Labels() []string {
	return ev.labels
}
