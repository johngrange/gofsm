package fsm

type fsmStateImpl struct {
	name        string
	transitions []Transition
	onEntry     Action
	onExit      Action
	stateLabels []string
	entryLabels []string
	exitLabels  []string
}

func NewState(name string, labels ...string) StateBuilder {
	s := &fsmStateImpl{
		name:        name,
		transitions: make([]Transition, 0),
		onEntry:     func(state State, fsmData interface{}, dispatcher Dispatcher) {},
		onExit:      func(state State, fsmData interface{}, dispatcher Dispatcher) {},
		stateLabels: []string{},
		entryLabels: []string{},
		exitLabels:  []string{},
	}

	for _, l := range labels {
		s.stateLabels = append(s.stateLabels, l)
	}
	return s
}

func (s *fsmStateImpl) Name() string {
	return s.name
}

func (s *fsmStateImpl) OnEntry(f Action, labels ...string) State {
	for _, l := range labels {
		s.entryLabels = append(s.entryLabels, l)
	}
	s.onEntry = f
	return s
}
func (s *fsmStateImpl) doEntry(fsm FSM) {
	s.onEntry(s, fsm.GetData(), fsm.GetDispatcher())
}
func (s *fsmStateImpl) OnExit(f Action, labels ...string) State {
	for _, l := range labels {
		s.exitLabels = append(s.exitLabels, l)
	}
	s.onExit = f
	return s
}
func (s *fsmStateImpl) doExit(fsm FSM) {
	s.onExit(s, fsm.GetData(), fsm.GetDispatcher())
}

func (s *fsmStateImpl) AddTransition(target State, labels ...string) TransitionBuilder {
	t := newTransition(s, target, labels...)

	s.transitions = append(s.transitions, t)
	return t
}
func (s *fsmStateImpl) Transitions() []Transition {
	return s.transitions
}

func (s *fsmStateImpl) StateLabels() []string {
	return s.stateLabels
}
func (s *fsmStateImpl) EntryLabels() []string {
	return s.entryLabels
}
func (s *fsmStateImpl) ExitLabels() []string {
	return s.exitLabels
}
