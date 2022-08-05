package fsm

type immediateFSMImpl struct {
	running      bool
	states       []FSMState
	currentState FSMState
	fsmData      interface{}
	tracers      []Tracer
	doTrace      bool
}

func NewImmediateFSM(initialState FSMState, data interface{}) ImmediateFSM {

	return &immediateFSMImpl{
		running:      false,
		states:       []FSMState{initialState},
		currentState: initialState,
		fsmData:      data,
		tracers:      make([]Tracer, 0),
	}
}

func (f *immediateFSMImpl) AddTracer(t Tracer) FSM {
	f.tracers = append(f.tracers, t)
	return f
}

func (f *immediateFSMImpl) traceTransition(ev Event, source, target FSMState) {
	for _, t := range f.tracers {
		t.OnTransition(ev, source, target, f.fsmData)
	}
}
func (f *immediateFSMImpl) AddState(s FSMState) FSM {
	f.states = append(f.states, s)
	return f
}

func (f *immediateFSMImpl) Start() {
	f.running = true
	f.traceOnEntry(f.currentState, f.fsmData)
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
	f.currentState.doExit(f.fsmData)
	f.traceOnExit(f.currentState, f.fsmData)
	f.traceTransition(ev, f.currentState, transition.Target())
	f.currentState = transition.Target()
	f.currentState.doEntry(f.fsmData)
	f.traceOnEntry(f.currentState, f.fsmData)
}
func (f *immediateFSMImpl) CurrentState() FSMState {
	return f.currentState
}
func (f *immediateFSMImpl) Dispatch(ev Event) {
	if f.running {
		f.processEvent(ev)
	}
}

func (f *immediateFSMImpl) traceOnEntry(state FSMState, fsmData interface{}) {
	for _, t := range f.tracers {
		t.OnEntry(state, fsmData)
	}
}
func (f *immediateFSMImpl) traceOnExit(state FSMState, fsmData interface{}) {
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
