package fsm

type transitionImpl struct {
	source       State
	target       State
	guard        TransitionGuard
	action       TransitionEffect
	triggerEvent string
}

func newTransition(source, target State) TransitionBuilder {

	return &transitionImpl{
		source: source,
		target: target,
		guard: func(fsmData, eventData interface{}) bool {
			return true
		},
		action: func(ev Event, fsmData interface{}) {},
	}
}

func (t *transitionImpl) Source() State {
	return t.source
}
func (t *transitionImpl) Target() State {
	return t.target
}
func (t *transitionImpl) SetTrigger(eventName string) TransitionBuilder {
	t.triggerEvent = eventName
	return t
}

func (t *transitionImpl) SetGuard(guard TransitionGuard) TransitionBuilder {
	t.guard = guard
	return t
}
func (t *transitionImpl) SetEffect(action TransitionEffect) TransitionBuilder {
	t.action = action
	return t
}

func (t *transitionImpl) shouldTransitionEv(ev Event, fsmData interface{}) bool {
	return ev.Name() == t.triggerEvent && t.guard(fsmData, ev.Data())
}
func (t *transitionImpl) shouldTransitionNoEv(fsmData interface{}) bool {
	return t.triggerEvent == "" && t.guard(fsmData, nil)
}

func (t *transitionImpl) GetEventName() string {
	return t.triggerEvent
}

func (t *transitionImpl) doAction(ev Event, fsmData interface{}) {
	t.action(ev, fsmData)
}

func (t *transitionImpl) IsLocal() bool {
	return t.source == t.target
}
