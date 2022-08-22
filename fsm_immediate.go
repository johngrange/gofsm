package fsm

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"
)

type immediateFSMImpl struct {
	running              bool
	initialState         State // always populated
	finalState           State // may be nil
	states               []State
	currentState         State
	fsmData              interface{}
	tracers              []Tracer
	eventProcesingActive bool
	eventQueue           chan Event
	dispatcher           Dispatcher
}

func (f *immediateFSMImpl) AddTracer(t Tracer) {
	f.tracers = append(f.tracers, t)
}

func (f *immediateFSMImpl) traceTransition(ev Event, source, target State) {
	for _, t := range f.tracers {
		t.OnTransition(ev, source, target, f.fsmData)
	}
}

func (f *immediateFSMImpl) Start() {
	f.running = true
	f.traceOnEntry(f.currentState, f.fsmData)
	f.currentState.doEntry(f)
	f.runToWaitCondition()
}
func (f *immediateFSMImpl) Stop() {
	f.running = false // stop accepting events on queue
}

func (f *immediateFSMImpl) Tick() {
	if f.running {
		f.runToWaitCondition()
	}
}
func (f *immediateFSMImpl) runToWaitCondition() {
	// keep evaluating no event transitions until we can't exit the current state
	for {
		transitioned := false
		for _, transition := range f.currentState.Transitions() {
			if transition.shouldTransitionNoEv(f.fsmData) {
				f.doTransition(nil, transition)
				transitioned = true
				break
			}
		}
		if !transitioned {
			return
		}
	}
}
func (f *immediateFSMImpl) doTransition(ev Event, transition Transition) {

	// UML spec 14.2.3.4.5, 14.2.3.4.6
	// state is exited after exit action completes

	// we are in new state before transition effect and new state entry actions called

	// if local transition, do not call exit or entry actions

	oldState := f.currentState
	nextState := transition.Target()

	transition.doAction(ev, f)
	f.traceTransition(ev, f.currentState, transition.Target())

	if !transition.IsLocal() {
		oldState.doExit(f)
		f.traceOnExit(oldState, f.fsmData)
		f.currentState = nextState
		nextState.doEntry(f)
		f.traceOnEntry(nextState, f.fsmData)
	}
}
func (f *immediateFSMImpl) CurrentState() State {
	return f.currentState
}
func (f *immediateFSMImpl) Dispatch(ev Event) {
	if f.running {
		f.eventQueue <- ev
		f.processImmediateEventQueue()
	}
}
func (f *immediateFSMImpl) processImmediateEventQueue() {
	// Don't allow nested calls to this method.  If more events get dispatched
	// as a result of this processing, they will be added to the event queue and we
	// will deal with them at this top level.
	if f.eventProcesingActive {
		return
	}
	f.eventProcesingActive = true
	defer func() {
		f.eventProcesingActive = false
	}()
	for len(f.eventQueue) > 0 {
		fmt.Fprintf(ginkgo.GinkgoWriter, "evq %d\n", len(f.eventQueue))
		ev := <-f.eventQueue
		f.processEvent(ev)
	}
}
func (f *immediateFSMImpl) traceOnEntry(state State, fsmData interface{}) {
	for _, t := range f.tracers {
		t.OnEntry(state, fsmData)
	}
}
func (f *immediateFSMImpl) traceOnExit(state State, fsmData interface{}) {
	for _, t := range f.tracers {
		t.OnExit(state, fsmData)
	}
}

func (f *immediateFSMImpl) traceRejectedEvent(ev Event, state State, fsmData interface{}) {
	for _, t := range f.tracers {
		t.OnRejectedEvent(ev, state, fsmData)
	}
}

func (f *immediateFSMImpl) processEvent(ev Event) {
	for _, transition := range f.currentState.Transitions() {
		if transition.shouldTransitionEv(ev, f.fsmData) {
			f.doTransition(ev, transition)
			f.runToWaitCondition()
			return
		}
	}
	f.traceRejectedEvent(ev, f.currentState, f.fsmData)
}

func (f *immediateFSMImpl) Visit(v Visitor) {
	for _, state := range f.states {
		v.VisitState(state)
		for _, transition := range state.Transitions() {
			v.VisitTransition(transition)
		}
	}
}

func (f *immediateFSMImpl) GetData() interface{} {
	return f.fsmData
}

func (f *immediateFSMImpl) GetDispatcher() Dispatcher {
	return f.dispatcher
}
