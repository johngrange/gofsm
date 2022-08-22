package fsm

type transitionBuilderImpl struct {
	source              StateBuilder
	target              StateBuilder
	guard               TransitionGuard
	action              TransitionEffect
	triggerEvent        string
	labels              []string
	triggerLabels       []string
	guardLabels         []string
	effectLabels        []string
	finalisedTransition Transition
}

func newTransitionBuilder(sourceStateBuilder, targetStateBuilder StateBuilder, labels ...string) TransitionBuilder {
	tb := &transitionBuilderImpl{
		source: sourceStateBuilder,
		target: targetStateBuilder,
		guard: func(fsmData, eventData interface{}) bool {
			return true
		},
		action:        func(ev Event, fsmData interface{}, dispatcher Dispatcher) {},
		labels:        []string{},
		triggerLabels: []string{},
		guardLabels:   []string{},
		effectLabels:  []string{},
	}
	for _, l := range labels {
		tb.labels = append(labels, l)
	}
	return tb
}

func (tb *transitionBuilderImpl) SetTrigger(eventName string, labels ...string) TransitionBuilder {
	for _, l := range labels {
		tb.triggerLabels = append(tb.triggerLabels, l)
	}
	tb.triggerEvent = eventName

	return tb
}
func (tb *transitionBuilderImpl) SetGuard(guard TransitionGuard, labels ...string) TransitionBuilder {
	for _, l := range labels {
		tb.guardLabels = append(tb.guardLabels, l)
	}
	tb.guard = guard
	return tb
}
func (tb *transitionBuilderImpl) SetEffect(effect TransitionEffect, labels ...string) TransitionBuilder {
	for _, l := range labels {
		tb.effectLabels = append(tb.effectLabels, l)
	}
	tb.action = effect

	return tb
}

func (tb *transitionBuilderImpl) Source() StateBuilder {
	return tb.source
}
func (tb *transitionBuilderImpl) Target() StateBuilder {
	return tb.target
}

func (tb *transitionBuilderImpl) build(source, target State) (Transition, error) {
	if tb.finalisedTransition != nil {
		return tb.finalisedTransition, nil
	}
	tb.finalisedTransition = &transitionImpl{
		source:        source,
		target:        target,
		guard:         tb.guard,
		action:        tb.action,
		triggerEvent:  tb.triggerEvent,
		labels:        tb.labels,
		triggerLabels: tb.triggerLabels,
		guardLabels:   tb.guardLabels,
		effectLabels:  tb.effectLabels,
	}
	return tb.finalisedTransition, nil
}
