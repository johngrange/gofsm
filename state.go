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

func (s *fsmStateImpl) Name() string {
	return s.name
}

func (s *fsmStateImpl) doExit(fsm FSM) {
	s.onExit(s, fsm.GetData(), fsm.GetDispatcher())
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

func (s *fsmStateImpl) doEntry(fsm FSM) {
	s.onEntry(s, fsm.GetData(), fsm.GetDispatcher())
}
