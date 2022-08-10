package fsm

type transitionImpl struct {
	source       FSMState
	target       FSMState
	guard        TransitionGuard
	action       TransitionAction
	triggerEvent string
}

func newTransition(source, target FSMState) Transition {

	return &transitionImpl{
		source: source,
		target: target,
		guard: func(state FSMState, fsmData, eventData interface{}) bool {
			return true
		},
		action: func(fromState, toState FSMState, ev Event, fsmData interface{}) {},
	}
}

func (t *transitionImpl) Source() FSMState {
	return t.source
}
func (t *transitionImpl) Target() FSMState {
	return t.target
}
func (t *transitionImpl) SetTrigger(eventName string) Transition {
	t.triggerEvent = eventName
	return t
}

func (t *transitionImpl) SetGuard(guard TransitionGuard) Transition {
	t.guard = guard
	return t
}
func (t *transitionImpl) SetAction(action TransitionAction) Transition {
	t.action = action
	return t
}

func (t *transitionImpl) shouldTransitionEv(ev Event, fsmData interface{}) bool {
	return ev.Name() == t.triggerEvent && t.guard(t.source, fsmData, ev.Data())
}
func (t *transitionImpl) shouldTransitionNoEv(fsmData interface{}) bool {
	return t.triggerEvent == "" && t.guard(t.source, fsmData, nil)
}

func (t *transitionImpl) GetEventName() string {
	return t.triggerEvent
}

func (t *transitionImpl) doAction(ev Event, fsmData interface{}) {
	t.action(t.source, t.target, ev, fsmData)
}
