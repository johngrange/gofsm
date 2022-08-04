package fsm

import (
	"fmt"
	"sync"
	"time"

	"github.com/onsi/ginkgo/v2"
)

type threadedFsmImpl struct {
	base       *immediateFSMImpl
	evalMX     sync.Mutex
	stop       chan struct{}
	eventQueue chan Event
}

const eventQueueLength = 50
const dataPollPeriod = time.Millisecond * 10

func NewThreadedFSM(initialState FSMState, data interface{}) FSM {
	return &threadedFsmImpl{
		base:       NewImmediateFSM(initialState, data).(*immediateFSMImpl),
		evalMX:     sync.Mutex{},
		eventQueue: make(chan Event, eventQueueLength),
	}
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
			fmt.Fprintln(ginkgo.GinkgoWriter, "stopping")
			return
		case ev := <-f.eventQueue:
			fmt.Fprintf(ginkgo.GinkgoWriter, "process event %+v\n", ev)
			f.base.processEvent(ev)
		case <-time.After(dataPollPeriod):
			fmt.Fprintf(ginkgo.GinkgoWriter, "run to wait\n")
			f.base.runToWaitCondition()
		}
	}
}
func (f *threadedFsmImpl) Dispatch(ev Event) {
	if f.base.running {
		f.eventQueue <- ev
	}
}

func (f *threadedFsmImpl) AddTracer(t Tracer) FSM {
	f.base.AddTracer(t)
	return f
}
func (f *threadedFsmImpl) AddState(s FSMState) FSM {
	f.base.AddState(s)
	return f
}
func (f *threadedFsmImpl) CurrentState() FSMState {
	return f.base.CurrentState()
}
