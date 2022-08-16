package fsm

import (
	"sync"
	"time"
)

type threadedFsmImpl struct {
	base       *immediateFSMImpl
	evalMX     sync.Mutex
	stop       chan struct{}
	eventQueue chan Event
	mx         sync.Mutex
}

const eventQueueLength = 50
const dataPollPeriod = time.Millisecond * 10

func NewThreadedFSM(data interface{}) FSMBuilder {
	fsm := &threadedFsmImpl{
		base:       NewImmediateFSM(data).(*immediateFSMImpl),
		evalMX:     sync.Mutex{},
		eventQueue: make(chan Event, eventQueueLength),
	}
	fsm.base.dispatcher = fsm
	return fsm
}

func (f *threadedFsmImpl) GetInitialState() StateBuilder {
	return f.base.GetInitialState()
}

func (f *threadedFsmImpl) GetFinalState() StateBuilder {
	return f.base.GetFinalState()
}

func (f *threadedFsmImpl) AddFinalState() StateBuilder {
	return f.base.AddFinalState()
}

func (f *threadedFsmImpl) Start() {
	f.stop = make(chan struct{})
	f.evalMX.Lock()
	defer f.evalMX.Unlock()
	f.base.Start()
	go f.runEventQueue()
}
func (f *threadedFsmImpl) Stop() {
	f.base.Stop()
	close(f.stop)
}
func (f *threadedFsmImpl) runEventQueue() {
	for {
		select {
		case <-f.stop:
			return
		case ev := <-f.eventQueue:
			f.mx.Lock()
			f.base.processEvent(ev)
			f.mx.Unlock()
		case <-time.After(dataPollPeriod):
			f.base.runToWaitCondition()
		}
	}
}
func (f *threadedFsmImpl) Dispatch(ev Event) {
	if f.base.running {
		f.eventQueue <- ev
	}
}

func (f *threadedFsmImpl) AddTracer(t Tracer) FSMBuilder {
	f.base.AddTracer(t)
	return f
}
func (f *threadedFsmImpl) AddState(s State) FSMBuilder {
	f.base.AddState(s)
	return f
}
func (f *threadedFsmImpl) CurrentState() State {
	return f.base.CurrentState()
}

func (f *threadedFsmImpl) Visit(v Visitor) {
	f.base.Visit(v)
}

func (f *threadedFsmImpl) GetData() interface{} {
	return f.base.fsmData
}

func (f *threadedFsmImpl) GetDispatcher() Dispatcher {
	return f
}
