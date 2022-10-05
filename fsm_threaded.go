package fsm

import (
	"fmt"
	"sync"
	"time"

	"github.com/onsi/ginkgo/v2"
)

type threadedFsmImpl struct {
	base                *immediateFSMImpl
	stop                chan struct{}
	eventQueue          chan Event
	mx, currStateMX     sync.RWMutex
	evaluateFSMChan     chan struct{} // entries in here trigger a re-evaluation of the FSM
	haltStateGoRoutines chan struct{} // closed when in-state go routines should exit (state change occurring).
	currentState        State
	currentStateChan    chan State
}

const eventQueueLength = 50
const dataPollPeriod = time.Millisecond * 10

func newThreadedFSM(base *immediateFSMImpl) FSM {
	fsm := &threadedFsmImpl{
		base:                base,
		eventQueue:          make(chan Event, eventQueueLength),
		evaluateFSMChan:     make(chan struct{}, eventQueueLength),
		haltStateGoRoutines: make(chan struct{}),
		currentStateChan:    make(chan State),
	}

	fsm.base.houseKeepStateEntry = func() {
		fsm.startTransitionTimers()
	}
	fsm.base.houseKeepStateExit = func() {
		fsm.stopTransitionTimers()
	}

	fsm.base.dispatcher = fsm
	fsm.currentState = fsm.base.currentState
	return fsm
}

func (f *threadedFsmImpl) Start() {
	f.stop = make(chan struct{})
	f.currStateMX.Lock()
	defer f.currStateMX.Unlock()
	f.base.Start()
	f.currentState = f.base.currentState
	go f.runEventQueue()
	go f.runCurrentStateChan()
}
func (f *threadedFsmImpl) Stop() {
	f.base.Stop()
	close(f.stop)
}

func (f *threadedFsmImpl) startTransitionTimers() {
	f.haltStateGoRoutines = make(chan struct{})
	for _, transition := range f.base.CurrentState().Transitions() {
		if transition.TriggerType() == TimerTrigger {
			go func() {
				select {
				case <-time.After(transition.TimerDuration()):
					fmt.Fprintf(ginkgo.GinkgoWriter, "%v timer expired for transition %s to %s\n", transition.TimerDuration(), transition.Source().Name(), transition.Target().Name())
					f.evaluateFSMChan <- struct{}{}
				case <-f.haltStateGoRoutines:
					fmt.Fprintf(ginkgo.GinkgoWriter, "%v timer cancelled for transition %s to %s\n", transition.TimerDuration(), transition.Source().Name(), transition.Target().Name())
					return
				}
			}()
		}
	}
}
func (f *threadedFsmImpl) stopTransitionTimers() {
	close(f.haltStateGoRoutines)
}

func (f *threadedFsmImpl) runCurrentStateChan() {
	for {
		select {
		case <-f.stop:
			return
		case s := <-f.currentStateChan:
			f.currStateMX.Lock()
			f.currentState = s
			f.currStateMX.Unlock()
		}
	}
}

func (f *threadedFsmImpl) runEventQueue() {
	for {
		select {
		case <-f.stop:
			return
		case ev := <-f.eventQueue:
			fmt.Fprintf(ginkgo.GinkgoWriter, "processing event %+v\n", ev)
			f.mx.Lock()
			fmt.Fprintf(ginkgo.GinkgoWriter, "current state before %+v\n", f.base.currentState)
			initialState := f.base.currentState
			f.base.processEvent(ev)
			fmt.Fprintf(ginkgo.GinkgoWriter, "current state after %+v\n", f.base.currentState)
			if initialState != f.base.currentState {
				f.currentStateChan <- f.base.currentState
			}
			f.mx.Unlock()
		case <-time.After(dataPollPeriod):
			f.evaluateFSMChan <- struct{}{}
		case <-f.evaluateFSMChan:
			// received instruction to re-evaluate FSM, so do so
			f.mx.Lock()
			initialState := f.base.currentState
			f.base.runToWaitCondition()
			if initialState != f.base.currentState {
				f.currentStateChan <- f.base.currentState
			}
			f.mx.Unlock()
		}
	}
}
func (f *threadedFsmImpl) Dispatch(ev Event) {
	if f.base.running {
		f.eventQueue <- ev
	}
}

func (f *threadedFsmImpl) AddTracer(t Tracer) {
	f.mx.Lock()
	f.base.AddTracer(t)
	f.mx.Unlock()
}

func (f *threadedFsmImpl) CurrentState() State {
	f.currStateMX.RLock()
	s := f.currentState
	f.currStateMX.RUnlock()
	return s
}

func (f *threadedFsmImpl) Visit(v Visitor) {
	f.mx.RLock()
	defer f.mx.RUnlock()
	f.base.Visit(v)
}

func (f *threadedFsmImpl) GetData() interface{} {
	// data object does not get changed once sm built
	return f.base.fsmData
}

func (f *threadedFsmImpl) GetDispatcher() Dispatcher {
	// does not get changed once sm built
	return f
}
