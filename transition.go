package fsm

import "time"

type transitionImpl struct {
	source         State
	target         State
	guard          TransitionGuard
	action         TransitionEffect
	triggerEvent   string
	labels         []string
	triggerLabels  []string
	guardLabels    []string
	effectLabels   []string
	triggerType    TriggerType
	timeoutTrigger time.Duration
	timerDeadline  time.Time
}

func (t *transitionImpl) Source() State {
	return t.source
}
func (t *transitionImpl) Target() State {
	return t.target
}

func (t *transitionImpl) shouldTransitionEv(ev Event, fsmData interface{}) bool {
	if t.triggerType != EventTrigger {
		return false
	}
	return ev.Name() == t.triggerEvent && t.guard(fsmData, ev.Data())
}
func (t *transitionImpl) shouldTransitionNoEv(fsmData interface{}) bool {
	switch t.triggerType {
	case EventTrigger:
		return false
	case NoTrigger:
		return t.guard(fsmData, nil)
	case TimerTrigger:
		return time.Now().After(t.timerDeadline) && t.guard(fsmData, nil)
	default:
		// shouldn't happen
		return false
	}
}

func (t *transitionImpl) startTimer(timeFrom time.Time) {
	if t.triggerType == TimerTrigger {
		t.timerDeadline = timeFrom.Add(t.timeoutTrigger)
	}
}
func (t *transitionImpl) TriggerType() TriggerType {
	return t.triggerType
}
func (t *transitionImpl) TimerDuration() time.Duration {
	if t.triggerType == TimerTrigger {
		return t.timeoutTrigger
	}
	return 0
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
