package fsm

import "errors"

type fsmBuilder struct {
	initialState       StateBuilder // always populated
	finalState         StateBuilder // may be nil
	stateBuilders      []StateBuilder
	fsmData            interface{}
	tracers            []Tracer
	finalisedImmediate ImmediateFSM
	finalisedThreaded  FSM
}

func NewFSMBuilder() FSMBuilder {
	initialState := NewStateBuilder(InitialStateName)
	return &fsmBuilder{

		initialState:  initialState,
		stateBuilders: []StateBuilder{},
		fsmData:       nil,
		tracers:       make([]Tracer, 0),
	}
}

func (b *fsmBuilder) SetData(data interface{}) FSMBuilder {
	b.fsmData = data
	return b
}
func (b *fsmBuilder) GetInitialState() StateBuilder {
	return b.initialState
}

func (b *fsmBuilder) GetFinalState() StateBuilder {
	return b.finalState
}

func (b *fsmBuilder) AddFinalState() StateBuilder {
	b.finalState = NewStateBuilder(FinalStateName)
	return b.finalState
}

func (b *fsmBuilder) BuildImmediateFSM() (ImmediateFSM, error) {
	if b.finalisedThreaded != nil {
		return nil, errors.New("builder already finalised as threaded fsm")
	}
	if b.finalisedImmediate != nil {
		return b.finalisedImmediate, nil
	}
	var err error
	b.finalisedImmediate, err = b.newImmediateFSMImpl()
	if err != nil {
		return nil, err
	}
	return b.finalisedImmediate, nil
}

func (b *fsmBuilder) newImmediateFSMImpl() (*immediateFSMImpl, error) {
	initialState, err := b.initialState.build()
	if err != nil {
		return nil, err
	}
	var finalState State
	if b.finalState != nil {
		finalState, err = b.finalState.build()
		if err != nil {
			return nil, err
		}
	}
	fsm := &immediateFSMImpl{
		initialState: initialState,
		running:      false,
		states:       []State{initialState},
		currentState: initialState,
		finalState:   finalState,
		fsmData:      b.fsmData,
		tracers:      b.tracers,
		eventQueue:   make(chan Event, eventQueueLength),
	}
	for _, stateBuilder := range b.stateBuilders {
		state, err := stateBuilder.build()
		if err != nil {
			return nil, err
		}
		fsm.states = append(fsm.states, state)
	}
	if finalState != nil {
		fsm.states = append(fsm.states, finalState)
	}

	// Build the transitions - they need concrete states to build
	// hence using two stage
	err = b.initialState.buildTransitions()
	if err != nil {
		return nil, err
	}
	if b.finalState != nil {
		err = b.finalState.buildTransitions()
		if err != nil {
			return nil, err
		}
	}
	for _, stateBuilder := range b.stateBuilders {
		err := stateBuilder.buildTransitions()
		if err != nil {
			return nil, err
		}
	}

	fsm.dispatcher = fsm
	return fsm, nil

}
func (b *fsmBuilder) BuildThreadedFSM() (FSM, error) {
	if b.finalisedImmediate != nil {
		return nil, errors.New("builder already finalised as immediate fsm")
	}
	if b.finalisedThreaded != nil {
		return b.finalisedThreaded, nil
	}
	imm, err := b.newImmediateFSMImpl()
	if err != nil {
		return nil, err
	}
	b.finalisedThreaded = newThreadedFSM(imm)
	return b.finalisedThreaded, nil
}

func (b *fsmBuilder) AddState(sb StateBuilder) FSMBuilder {
	b.stateBuilders = append(b.stateBuilders, sb)
	return b

}
func (b *fsmBuilder) NewState(name string, labels ...string) StateBuilder {
	sb := NewStateBuilder(name, labels...)
	b.stateBuilders = append(b.stateBuilders, sb)
	return sb
}

func (b *fsmBuilder) AddTracer(t Tracer) FSMBuilder {
	b.tracers = append(b.tracers, t)
	return b
}
