package fsm

import (
	"fmt"
	"io"
)

func RenderPlantUML(w io.Writer, stateMachine FSM) error {
	visitor := plantUMLVisitor{
		w:    w,
		errs: []error{},
	}
	_, err := fmt.Fprintln(w, "@startuml")
	if err != nil {
		return err
	}
	stateMachine.Visit(&visitor)
	if len(visitor.errs) > 0 {
		return visitor.errs[0]
	}
	_, err = fmt.Fprintln(w, "@enduml")
	if err != nil {
		return err
	}
	return nil
}

type plantUMLVisitor struct {
	w              io.Writer
	errs           []error
	seenFirstState bool
}

func (p *plantUMLVisitor) VisitState(state FSMState) {
	if !p.seenFirstState {
		fmt.Fprintf(p.w, "[*] --> %s\n", state.Name())
		p.seenFirstState = true
	}
}
func (p *plantUMLVisitor) VisitTransition(t Transition) {

	evName := t.GetEventName()
	if evName != "" {
		evName = " : " + evName
	}
	_, err := fmt.Fprintf(p.w, "%s --> %s%s\n", t.Source().Name(), t.Target().Name(), evName)
	if err != nil {
		p.errs = append(p.errs, err)
	}
}
