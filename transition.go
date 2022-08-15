package fsm

type transitionImpl struct {
	source        State
	target        State
	guard         TransitionGuard
	action        TransitionEffect
	triggerEvent  string
	labels        []string
	triggerLabels []string
	guardLabels   []string
	effectLabels  []string
}

func newTransition(source, target State, labels ...string) *transitionImpl {

	t := &transitionImpl{
		source: source,
		target: target,
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
		t.labels = append(labels, l)
	}
	return t
}

func (t *transitionImpl) Source() State {
	return t.source
}
func (t *transitionImpl) Target() State {
	return t.target
}
func (t *transitionImpl) SetTrigger(eventName string, labels ...string) TransitionBuilder {
	for _, l := range labels {
		t.triggerLabels = append(t.triggerLabels, l)
	}
	t.triggerEvent = eventName
	return t
}

func (t *transitionImpl) SetGuard(guard TransitionGuard, labels ...string) TransitionBuilder {
	for _, l := range labels {
		t.guardLabels = append(t.guardLabels, l)
	}
	t.guard = guard
	return t
}
func (t *transitionImpl) SetEffect(action TransitionEffect, labels ...string) TransitionBuilder {
	for _, l := range labels {
		t.effectLabels = append(t.effectLabels, l)
	}
	t.action = action
	return t
}

func (t *transitionImpl) shouldTransitionEv(ev Event, fsmData interface{}) bool {
	return ev.Name() == t.triggerEvent && t.guard(fsmData, ev.Data())
}
func (t *transitionImpl) shouldTransitionNoEv(fsmData interface{}) bool {
	return t.triggerEvent == "" && t.guard(fsmData, nil)
}

func (t *transitionImpl) EventName() string {
	return t.triggerEvent
}

func (t *transitionImpl) doAction(ev Event, fsm FSM) {
	t.action(ev, fsm.GetData(), fsm.GetDispatcher())
}

func (t *transitionImpl) IsLocal() bool {
	return t.source == t.target
}

func (t *transitionImpl) TriggerLabels() []string {
	return t.triggerLabels
}
func (t *transitionImpl) GuardLabels() []string {
	return t.guardLabels
}
func (t *transitionImpl) EffectLabels() []string {
	return t.effectLabels
}
func (t *transitionImpl) Labels() []string {
	return t.labels
}
