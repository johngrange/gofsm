package fsm

type fsmStateBuilder struct {
	name           string
	transitions    []TransitionBuilder
	onEntry        Action
	onExit         Action
	stateLabels    []string
	entryLabels    []string
	exitLabels     []string
	finalisedState *fsmStateImpl
}

func NewStateBuilder(name string, labels ...string) StateBuilder {
	sb := &fsmStateBuilder{
		name:        name,
		transitions: make([]TransitionBuilder, 0),
		onEntry:     func(state State, fsmData interface{}, dispatcher Dispatcher) {},
		onExit:      func(state State, fsmData interface{}, dispatcher Dispatcher) {},
		stateLabels: []string{},
		entryLabels: []string{},
		exitLabels:  []string{},
	}

	sb.stateLabels = append(sb.stateLabels, labels...)
	return sb
}

func (sb *fsmStateBuilder) OnEntry(f Action, labels ...string) StateBuilder {
	sb.entryLabels = append(sb.entryLabels, labels...)
	sb.onEntry = f
	return sb
}
func (sb *fsmStateBuilder) OnExit(f Action, labels ...string) StateBuilder {
	sb.exitLabels = append(sb.exitLabels, labels...)
	sb.onExit = f
	return sb
}

func (sb *fsmStateBuilder) AddTransition(target StateBuilder, labels ...string) TransitionBuilder {
	t := newTransitionBuilder(sb, target, labels...)

	sb.transitions = append(sb.transitions, t)
	return t
}

func (sb *fsmStateBuilder) build() (State, error) {
	if sb.finalisedState != nil {
		return sb.finalisedState, nil
	}
	state := &fsmStateImpl{
		name:        sb.name,
		transitions: make([]Transition, 0),
		onEntry:     sb.onEntry,
		onExit:      sb.onExit,
		stateLabels: sb.stateLabels,
		entryLabels: sb.entryLabels,
		exitLabels:  sb.exitLabels,
	}
	sb.finalisedState = state
	return state, nil
}

func (sb *fsmStateBuilder) buildTransitions() error {
	for _, tb := range sb.transitions {
		var source, target State
		var err error
		source, err = tb.Source().build()
		if err != nil {
			return err
		}
		target, err = tb.Target().build()
		if err != nil {
			return err
		}
		transition, err := tb.build(source, target)
		if err != nil {
			return err
		}
		sb.finalisedState.transitions = append(sb.finalisedState.transitions, transition)
	}

	return nil
}
