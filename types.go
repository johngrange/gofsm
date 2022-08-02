package fsm

import "time"

type FSMTraceEntry struct {
	TransitionTime           time.Time
	EventName                string
	SourceState, TargetState string
}

type FSM interface {
	AddState(FSMState) FSM
	Dispatch(Event)
	CurrentState() FSMState
	Start()
	Stop()
	Trace(bool) FSM
	GetTrace() []FSMTraceEntry
}

type Event interface {
	Name() string
	Data() interface{}
}

type FSMState interface {
	AddTransition(target FSMState) Transition
	OnEntry(StateEntryFunc) FSMState
	OnExit(StateExitFunc) FSMState
	Name() string
	GetTransitions() []Transition
	doExit(fsmData interface{})
	doEntry(fsmData interface{})
}

type StateEntryFunc func(state FSMState, fsmData interface{})
type StateExitFunc func(state FSMState, fsmData interface{})

type TransitionGuard func(fsmData, eventData interface{}) bool

type Transition interface {
	Source() FSMState
	Target() FSMState
	SetTrigger(eventName string) Transition
	SetGuard(TransitionGuard) Transition
	shouldTransitionEv(ev Event, fsmData interface{}) bool // If this transition accepts supplied event and guard is met, then return true
	shouldTransitionNoEv(fsmData interface{}) bool         // If this transition guard is met, with no need for event, then return true.  will always return false if trigger event set.
}
