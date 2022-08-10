package fsm

type fsmStateImpl struct {
	name        string
	transitions []Transition
	onEntry     Action
	onExit      Action
}

func NewState(name string) FSMState {
	return &fsmStateImpl{
		name:        name,
		transitions: make([]Transition, 0),
		onEntry:     func(state FSMState, fsmData interface{}) {},
		onExit:      func(state FSMState, fsmData interface{}) {},
	}
}

func (s *fsmStateImpl) Name() string {
	return s.name
}

func (s *fsmStateImpl) OnEntry(f Action) FSMState {
	s.onEntry = f
	return s
}
func (s *fsmStateImpl) doEntry(fsmData interface{}) {
	s.onEntry(s, fsmData)
}
func (s *fsmStateImpl) OnExit(f Action) FSMState {
	s.onExit = f
	return s
}
func (s *fsmStateImpl) doExit(fsmData interface{}) {
	s.onExit(s, fsmData)
}

func (s *fsmStateImpl) AddTransition(target FSMState) Transition {
	t := newTransition(s, target)
	s.transitions = append(s.transitions, t)
	return t
}
func (s *fsmStateImpl) GetTransitions() []Transition {
	return s.transitions
}
