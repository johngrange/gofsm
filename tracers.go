package fsm

import (
	"fmt"
	"time"
)

type StateCounter struct {
	Counts map[string]uint64
}

func NewStateCounter() *StateCounter {
	return &StateCounter{
		Counts: make(map[string]uint64),
	}
}

func (s *StateCounter) OnEntry(state FSMState, fsmData interface{}) {
	val := s.Counts[state.Name()]
	val++
	s.Counts[state.Name()] = val
}
func (s *StateCounter) OnExit(state FSMState, fsmData interface{}) {

}
func (s *StateCounter) OnTransition(ev Event, sourceState, targetState FSMState, fsmData interface{}) {

}

type FSMLogEntry struct {
	When    time.Time
	Message string
}

type FSMLogger struct {
	Entries []FSMLogEntry
}

func NewFSMLogger() *FSMLogger {
	return &FSMLogger{
		Entries: make([]FSMLogEntry, 0),
	}
}

func (l *FSMLogger) OnEntry(state FSMState, fsmData interface{}) {
	l.Entries = append(l.Entries, FSMLogEntry{
		time.Now(),
		fmt.Sprintf("Entered State: %s, state: %+v, fsm: %+v", state.Name(), state, fsmData),
	})
}
func (l *FSMLogger) OnExit(state FSMState, fsmData interface{}) {
	l.Entries = append(l.Entries, FSMLogEntry{
		time.Now(),
		fmt.Sprintf("Exited State: %s, state: %+v, fsm: %+v", state.Name(), state, fsmData),
	})
}
func (l *FSMLogger) OnTransition(ev Event, sourceState, targetState FSMState, fsmData interface{}) {
	l.Entries = append(l.Entries, FSMLogEntry{
		time.Now(),
		fmt.Sprintf("Transitioning event: %+v, Source, %s: %+v+ Target: %s:%+v, fsm: %+v", ev, sourceState.Name(), sourceState, targetState.Name(), targetState, fsmData),
	})

}
