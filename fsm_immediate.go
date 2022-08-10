package fsm

type immediateFSMImpl struct {
	running      bool
	states       []State
	currentState State
	fsmData      interface{}
	tracers      []Tracer
	doTrace      bool
}

func NewImmediateFSM(initialState State, data interface{}) ImmediateFSMBuilder {

	return &immediateFSMImpl{
		running:      false,
		states:       []State{initialState},
		currentState: initialState,
		fsmData:      data,
		tracers:      make([]Tracer, 0),
	}
}

func (f *immediateFSMImpl) AddTracer(t Tracer) ImmediateFSMBuilder {
	f.tracers = append(f.tracers, t)
	return f
}

func (f *immediateFSMImpl) traceTransition(ev Event, source, target State) {
	for _, t := range f.tracers {
		t.OnTransition(ev, source, target, f.fsmData)
	}
}
func (f *immediateFSMImpl) AddState(s State) ImmediateFSMBuilder {
	f.states = append(f.states, s)
	return f
}

func (f *immediateFSMImpl) Start() {
	f.running = true
	f.traceOnEntry(f.currentState, f.fsmData)
	f.currentState.doEntry(f.fsmData)
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
		for _, transition := range f.currentState.GetTransitions() {
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

	transition.doAction(ev, f.fsmData)
	f.traceTransition(ev, f.currentState, transition.Target())

	if !transition.IsLocal() {
		oldState.doExit(f.fsmData)
		f.traceOnExit(oldState, f.fsmData)
		f.currentState = nextState
		nextState.doEntry(f.fsmData)
		f.traceOnEntry(nextState, f.fsmData)
	}
}
func (f *immediateFSMImpl) CurrentState() State {
	return f.currentState
}
func (f *immediateFSMImpl) Dispatch(ev Event) {
	if f.running {
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

func (f *immediateFSMImpl) processEvent(ev Event) {
	for _, transition := range f.currentState.GetTransitions() {
		if transition.shouldTransitionEv(ev, f.fsmData) {
			f.doTransition(ev, transition)
			f.runToWaitCondition()
		}
	}
}

func (f *immediateFSMImpl) Visit(v Visitor) {
	for _, state := range f.states {
		v.VisitState(state)
		for _, transition := range state.GetTransitions() {
			v.VisitTransition(transition)
		}
	}
}

func (f *immediateFSMImpl) GetData() interface{} {
	return f.fsmData
}
