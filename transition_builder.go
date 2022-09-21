package fsm

import "time"

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
	triggerType         TriggerType
	timeoutTrigger      time.Duration
}

func newTransitionBuilder(sourceStateBuilder, targetStateBuilder StateBuilder, labels ...string) TransitionBuilder {
	tb := &transitionBuilderImpl{
		source: sourceStateBuilder,
		target: targetStateBuilder,
		guard: func(fsmData, eventData interface{}) bool {
			return true
		},
		action:         func(ev Event, fsmData interface{}, dispatcher Dispatcher) {},
		labels:         []string{},
		triggerLabels:  []string{},
		guardLabels:    []string{},
		effectLabels:   []string{},
		triggerType:    NoTrigger,
		timeoutTrigger: 0,
	}
	tb.labels = append(tb.labels, labels...)
	return tb
}

func (tb *transitionBuilderImpl) SetEventTrigger(eventName string, labels ...string) TransitionBuilder {
	tb.triggerLabels = append(tb.triggerLabels, labels...)
	tb.triggerEvent = eventName
	tb.triggerType = EventTrigger

	return tb
}

func (tb *transitionBuilderImpl) SetTimedTrigger(timer time.Duration, labels ...string) TransitionBuilder {
	tb.triggerLabels = append(tb.triggerLabels, labels...)
	tb.timeoutTrigger = timer
	tb.triggerType = TimerTrigger
	return tb
}
func (tb *transitionBuilderImpl) SetGuard(guard TransitionGuard, labels ...string) TransitionBuilder {
	tb.guardLabels = append(tb.guardLabels, labels...)
	tb.guard = guard
	return tb
}
func (tb *transitionBuilderImpl) SetEffect(effect TransitionEffect, labels ...string) TransitionBuilder {
	tb.effectLabels = append(tb.effectLabels, labels...)
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
		source:         source,
		target:         target,
		guard:          tb.guard,
		action:         tb.action,
		triggerEvent:   tb.triggerEvent,
		labels:         tb.labels,
		triggerLabels:  tb.triggerLabels,
		guardLabels:    tb.guardLabels,
		effectLabels:   tb.effectLabels,
		triggerType:    tb.triggerType,
		timeoutTrigger: tb.timeoutTrigger,
	}
	return tb.finalisedTransition, nil
}

func (tb *transitionBuilderImpl) TriggerType() TriggerType {
	return tb.triggerType
}
