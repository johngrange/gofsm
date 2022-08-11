package fsm

type fsmStateImpl struct {
	name        string
	transitions []Transition
	onEntry     Action
	onExit      Action
}

func NewState(name string) StateBuilder {
	return &fsmStateImpl{
		name:        name,
		transitions: make([]Transition, 0),
		onEntry:     func(state State, fsmData interface{}, dispatcher Dispatcher) {},
		onExit:      func(state State, fsmData interface{}, dispatcher Dispatcher) {},
	}
}

func (s *fsmStateImpl) Name() string {
	return s.name
}

func (s *fsmStateImpl) OnEntry(f Action) State {
	s.onEntry = f
	return s
}
func (s *fsmStateImpl) doEntry(fsm FSM) {
	s.onEntry(s, fsm.GetData(), fsm.GetDispatcher())
}
func (s *fsmStateImpl) OnExit(f Action) State {
	s.onExit = f
	return s
}
func (s *fsmStateImpl) doExit(fsm FSM) {
	s.onExit(s, fsm.GetData(), fsm.GetDispatcher())
}

func (s *fsmStateImpl) AddTransition(target State) TransitionBuilder {
	t := newTransition(s, target)
	s.transitions = append(s.transitions, t)
	return t
}
func (s *fsmStateImpl) GetTransitions() []Transition {
	return s.transitions
}
