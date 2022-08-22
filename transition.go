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

func (t *transitionImpl) Source() State {
	return t.source
}
func (t *transitionImpl) Target() State {
	return t.target
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
