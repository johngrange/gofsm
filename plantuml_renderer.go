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

func (p *plantUMLVisitor) VisitState(state State) {
	stateName := state.Name()
	if stateName == InitialStateName {
		stateName = "[*]"
	}

	for _, l := range state.StateLabels() {
		fmt.Fprintf(p.w, "%s : %s\n", stateName, l)
	}
	for _, l := range state.EntryLabels() {
		fmt.Fprintf(p.w, "%s : entry/%s\n", stateName, l)
	}
	for _, l := range state.ExitLabels() {
		fmt.Fprintf(p.w, "%s : exit/%s\n", stateName, l)
	}
}
func (p *plantUMLVisitor) VisitTransition(t Transition) {

	evName := t.EventName()
	if evName != "" {
		evName = " : " + evName
	}
	guard := ""
	if len(t.GuardLabels()) > 0 {
		guard = " "
		for _, l := range t.GuardLabels() {
			guard += fmt.Sprintf("[%s] ", l)
		}
	}
	effect := ""
	if len(t.EffectLabels()) > 0 {
		effect = "/"
		for idx, l := range t.EffectLabels() {
			effect += l
			if idx < len(t.EffectLabels())-1 {
				effect += " "
			}
		}

	}

	sourceName := t.Source().Name()
	if sourceName == InitialStateName {
		sourceName = "[*]"
	}
	targetName := t.Target().Name()
	if targetName == FinalStateName {
		targetName = "[*]"
	}
	_, err := fmt.Fprintf(p.w, "%s --> %s%s%s%s\n", sourceName, targetName, evName, guard, effect)
	if err != nil {
		p.errs = append(p.errs, err)
	}
}
