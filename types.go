package fsm

import (
	"time"
)

const (
	InitialStateName = "initial"
	FinalStateName   = "FinalState"
)

type FSMTraceEntry struct {
	TransitionTime           time.Time
	EventName                string
	SourceState, TargetState string
}

type FSMBuilder interface {
	AddState(StateBuilder) FSMBuilder
	NewState(name string, labels ...string) StateBuilder
	AddTracer(Tracer) FSMBuilder
	AddFinalState() StateBuilder
	GetInitialState() StateBuilder
	GetFinalState() StateBuilder
	BuildImmediateFSM() (ImmediateFSM, error)
	BuildThreadedFSM() (FSM, error)
	SetData(data interface{}) FSMBuilder
}

type Dispatcher interface {
	Dispatch(Event)
}

type FSM interface {
	Dispatcher
	Visitable
	Observable

	CurrentState() State
	Start()
	Stop()
	GetData() interface{}
	GetDispatcher() Dispatcher
}

type ImmediateFSM interface {
	FSM
	Tick() // Manually check for and progress state changes that are not event driven
}

type Event interface {
	Name() string
	Data() interface{}
	Labels() []string
}

type StateBuilder interface {
	AddTransition(target StateBuilder, labels ...string) TransitionBuilder
	OnEntry(action Action, labels ...string) StateBuilder
	OnExit(action Action, labels ...string) StateBuilder
	build() (State, error)
	buildTransitions() error
}

type State interface {
	Name() string
	Transitions() []Transition
	StateLabels() []string
	EntryLabels() []string
	ExitLabels() []string
	doExit(fsm FSM)
	doEntry(fsm FSM)
}

type Tracer interface {
	OnEntry(state State, fsmData interface{})
	OnExit(state State, fsmData interface{})
	OnTransition(ev Event, sourceState, targetState State, fsmData interface{})
	OnRejectedEvent(ev Event, state State, fmsData interface{})
}

type Action func(state State, fsmData interface{}, dispatcher Dispatcher)
type TransitionEffect func(ev Event, fsmData interface{}, dispatcher Dispatcher)
type TransitionGuard func(fsmData, eventData interface{}) bool

type TransitionBuilder interface {
	SetTrigger(eventName string, labels ...string) TransitionBuilder
	SetGuard(guard TransitionGuard, labels ...string) TransitionBuilder
	SetEffect(efffect TransitionEffect, labels ...string) TransitionBuilder
	Source() StateBuilder
	Target() StateBuilder
	build(source, target State) (Transition, error)
}

type Transition interface {
	Source() State
	Target() State
	IsLocal() bool
	TriggerLabels() []string
	GuardLabels() []string
	EffectLabels() []string
	EventName() string
	shouldTransitionEv(ev Event, fsmData interface{}) bool // If this transition accepts supplied event and guard is met, then return true
	shouldTransitionNoEv(fsmData interface{}) bool         // If this transition guard is met, with no need for event, then return true.  will always return false if trigger event set.
	doAction(ev Event, fsm FSM)
}

type Visitable interface {
	Visit(Visitor)
}
type Observable interface {
	AddTracer(Tracer)
}

type Visitor interface {
	VisitState(state State)
	VisitTransition(transition Transition)
}
