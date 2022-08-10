package fsm

import (
	"time"
)

type FSMTraceEntry struct {
	TransitionTime           time.Time
	EventName                string
	SourceState, TargetState string
}

type FSM interface {
	Visitable
	AddState(FSMState) FSM
	AddTracer(Tracer) FSM
	Dispatch(Event)
	CurrentState() FSMState
	Start()
	Stop()
	GetData() interface{}
}

type ImmediateFSM interface {
	FSM
	Tick() // Manually check for and progress state changes that are not event driven
}

type Event interface {
	Name() string
	Data() interface{}
}

type StateDataFactory func() interface{}

type FSMState interface {
	AddTransition(target FSMState) Transition
	OnEntry(Action) FSMState
	OnExit(Action) FSMState
	SetDataFactory(StateDataFactory)
	GetCurrentData() interface{}
	Name() string
	GetTransitions() []Transition
	doExit(fsmData interface{})
	doEntry(fsmData interface{})
	initialiseStateData()
}

func DefaultStateDataFactory() interface{} {
	return nil
}

type Tracer interface {
	OnEntry(state FSMState, fsmData interface{})
	OnExit(state FSMState, fsmData interface{})
	OnTransition(ev Event, sourceState, targetState FSMState, fsmData interface{})
}
type Action func(state FSMState, fsmData interface{})
type TransitionAction func(fromState, toState FSMState, ev Event, fsmData interface{})
type TransitionGuard func(fromState FSMState, fsmData, eventData interface{}) bool

type Transition interface {
	Source() FSMState
	Target() FSMState
	SetTrigger(eventName string) Transition
	SetGuard(TransitionGuard) Transition
	SetAction(TransitionAction) Transition
	GetEventName() string
	shouldTransitionEv(ev Event, fsmData interface{}) bool // If this transition accepts supplied event and guard is met, then return true
	shouldTransitionNoEv(fsmData interface{}) bool         // If this transition guard is met, with no need for event, then return true.  will always return false if trigger event set.
	doAction(ev Event, fsmData interface{})
}

type Visitable interface {
	Visit(Visitor)
}

type Visitor interface {
	VisitState(state FSMState)
	VisitTransition(transition Transition)
}
