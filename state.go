package fsm

type fsmStateImpl struct {
	name        string
	transitions []Transition
	onEntry     Action
	onExit      Action
	dataFactory StateDataFactory
	currentData interface{}
}

func NewState(name string) FSMState {
	return &fsmStateImpl{
		name:        name,
		transitions: make([]Transition, 0),
		onEntry:     func(state FSMState, fsmData interface{}) {},
		onExit:      func(state FSMState, fsmData interface{}) {},
		dataFactory: DefaultStateDataFactory,
		currentData: nil,
	}
}

func (s *fsmStateImpl) Name() string {
	return s.name
}

func (s *fsmStateImpl) SetDataFactory(f StateDataFactory) {
	s.dataFactory = f
}
func (s *fsmStateImpl) GetCurrentData() interface{} {
	return s.currentData
}

func (s *fsmStateImpl) OnEntry(f Action) FSMState {
	s.onEntry = f
	return s
}
func (s *fsmStateImpl) doEntry(fsmData interface{}) {
	s.onEntry(s, fsmData)
}

func (s *fsmStateImpl) initialiseStateData() {
	s.currentData = s.dataFactory()
}
func (s *fsmStateImpl) OnExit(f Action) FSMState {
	s.onExit = f
	return s
}
func (s *fsmStateImpl) doExit(fsmData interface{}) {
	s.onExit(s, fsmData)
	s.currentData = nil
}

func (s *fsmStateImpl) AddTransition(target FSMState) Transition {
	t := newTransition(s, target)
	s.transitions = append(s.transitions, t)
	return t
}
func (s *fsmStateImpl) GetTransitions() []Transition {
	return s.transitions
}
