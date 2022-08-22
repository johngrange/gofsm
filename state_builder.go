package fsm

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"
)

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

	for _, l := range labels {
		sb.stateLabels = append(sb.stateLabels, l)
	}
	return sb
}

func (sb *fsmStateBuilder) OnEntry(f Action, labels ...string) StateBuilder {
	for _, l := range labels {
		sb.entryLabels = append(sb.entryLabels, l)
	}
	sb.onEntry = f
	return sb
}
func (sb *fsmStateBuilder) OnExit(f Action, labels ...string) StateBuilder {
	for _, l := range labels {
		sb.exitLabels = append(sb.exitLabels, l)
	}
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
		fmt.Fprintf(ginkgo.GinkgoWriter, "returning built state %s, transitions: %d\n", sb.finalisedState.Name(), len(sb.finalisedState.Transitions()))
		return sb.finalisedState, nil
	}
	fmt.Fprintf(ginkgo.GinkgoWriter, "building state %s\n", sb.name)
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
	fmt.Fprintf(ginkgo.GinkgoWriter, "building transitions for state %s, with %d transitions\n", sb.name, len(sb.transitions))

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
		sb.finalisedState.transitions = append(sb.finalisedState.transitions, transition)
	}
	fmt.Fprintf(ginkgo.GinkgoWriter, "state %s has %d transitions after building\n", sb.finalisedState.name, len(sb.finalisedState.transitions))

	return nil
}
