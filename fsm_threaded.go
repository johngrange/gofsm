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

func NewThreadedFSM(initialState State, data interface{}) FSMBuilder {
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
