package fsm

import (
	"sync"
	"time"
)

type fsmImpl struct {
	running      bool
	states       []FSMState
	currentState FSMState
	fsmData      interface{}
	trace        []FSMTraceEntry
	doTrace      bool
	evalMX       sync.Mutex
	stop         chan struct{}
	eventQueue   chan Event
}

const eventQueueLength = 50
const dataPollPeriod = time.Millisecond * 10

func NewThreadedFSM(initialState FSMState, data interface{}) FSM {
	return &fsmImpl{
		running:      false,
		states:       make([]FSMState, 0),
		currentState: initialState,
		fsmData:      data,
		trace:        make([]FSMTraceEntry, 0),
		doTrace:      false,
		evalMX:       sync.Mutex{},
		eventQueue:   make(chan Event, eventQueueLength),
	}
}

func (f *fsmImpl) GetTrace() []FSMTraceEntry {
	return f.trace
}
func (f *fsmImpl) Trace(enable bool) FSM {
	f.doTrace = enable
	return f
}
func (f *fsmImpl) traceTransition(ev Event, source, target FSMState) {
	if f.doTrace {
		evName := ""
		if ev != nil {
			evName = ev.Name()
		}
		f.trace = append(f.trace, FSMTraceEntry{
			TransitionTime: time.Now(),
			EventName:      evName,
			SourceState:    source.Name(),
			TargetState:    target.Name(),
		})
	}
}
func (f *fsmImpl) AddState(s FSMState) FSM {
	f.states = append(f.states, s)
	return f
}

func (f *fsmImpl) Start() {
	f.running = true
	f.stop = make(chan struct{})
	f.evalMX.Lock()
	defer f.evalMX.Unlock()
	f.runToWaitCondition()
	go f.runEventQueue()
}
func (f *fsmImpl) Stop() {
	f.running = false // stop accepting events on queue
	close(f.stop)
}
func (f *fsmImpl) runEventQueue() {
	for {
		select {
		case <-f.stop:
			return
		case ev := <-f.eventQueue:
			f.processEvent(ev)
		case <-time.After(dataPollPeriod):
			f.evalMX.Lock()
			f.runToWaitCondition()
			f.evalMX.Unlock()
		}
	}
}
func (f *fsmImpl) runToWaitCondition() {
	// keep evaluating no event transitions until we can't exit the current state
	for {
		transitioned := false
		for _, transition := range f.currentState.GetTransitions() {
			if transition.shouldTransitionNoEv(f.fsmData) {
				f.currentState.doExit(f.fsmData)
				f.traceTransition(nil, f.currentState, transition.Target())
				f.currentState = transition.Target()
				f.currentState.doEntry(f.fsmData)
				transitioned = true
				break
			}
		}
		if !transitioned {
			return
		}
	}
}

func (f *fsmImpl) CurrentState() FSMState {
	return f.currentState
}
func (f *fsmImpl) Dispatch(ev Event) {
	if f.running {
		f.eventQueue <- ev
	}
}
func (f *fsmImpl) processEvent(ev Event) {
	f.evalMX.Lock()
	defer f.evalMX.Unlock()
	for _, transition := range f.currentState.GetTransitions() {
		if transition.shouldTransitionEv(ev, f.fsmData) {
			f.currentState.doExit(f.fsmData)
			f.traceTransition(ev, f.currentState, transition.Target())
			f.currentState = transition.Target()
			f.currentState.doEntry(f.fsmData)
			f.runToWaitCondition()
		}
	}
}
