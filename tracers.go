package fsm

import (
	"fmt"
	"time"
)

type StateCounter struct {
	StateCounts         map[string]uint64
	RejectedEventCounts map[string]uint64
}

func NewStateCounter() *StateCounter {
	return &StateCounter{
		StateCounts:         make(map[string]uint64),
		RejectedEventCounts: make(map[string]uint64),
	}
}

func (s *StateCounter) OnEntry(state State, fsmData interface{}) {
	val := s.StateCounts[state.Name()]
	val++
	s.StateCounts[state.Name()] = val
}
func (s *StateCounter) OnExit(state State, fsmData interface{}) {

}
func (s *StateCounter) OnTransition(ev Event, sourceState, targetState State, fsmData interface{}) {

}

func (s *StateCounter) OnRejectedEvent(ev Event, state State, fsmData interface{}) {
	count := s.RejectedEventCounts[ev.Name()]
	count++
	s.RejectedEventCounts[ev.Name()] = count
}

type LogEntry struct {
	When    time.Time
	Message string
}

type Logger struct {
	Entries  []LogEntry
	Detailed bool
}

func NewFSMLogger() *Logger {
	return &Logger{
		Entries: make([]LogEntry, 0),
	}
}

func (l *Logger) OnEntry(state State, fsmData interface{}) {
	detail := ""
	if l.Detailed {
		detail = fmt.Sprintf(": state: %+v, fsm: %+v", state, fsmData)
	}
	l.Entries = append(l.Entries, LogEntry{
		time.Now(),
		fmt.Sprintf("En  : %s%s", state.Name(), detail),
	})
}
func (l *Logger) OnExit(state State, fsmData interface{}) {
	detail := ""
	if l.Detailed {
		detail = fmt.Sprintf(": state: %+v, fsm: %+v", state, fsmData)
	}
	l.Entries = append(l.Entries, LogEntry{
		time.Now(),
		fmt.Sprintf("Ex  : %s%s", state.Name(), detail),
	})
}
func (l *Logger) OnTransition(ev Event, sourceState, targetState State, fsmData interface{}) {
	detail := ""
	if l.Detailed {
		detail = fmt.Sprintf(":  source, %+v, target: %+v, fsm: %+v", sourceState, targetState, fsmData)
	}
	evName := ""
	if ev != nil {
		evName = ev.Name()
	}
	l.Entries = append(l.Entries, LogEntry{
		time.Now(),
		fmt.Sprintf("Tr  : %s --> %s [%s]%s", sourceState.Name(), targetState.Name(), evName, detail),
	})
}

func (l *Logger) OnRejectedEvent(ev Event, state State, fsmData interface{}) {
	detail := ""
	if l.Detailed {
		detail = fmt.Sprintf(":  event, %+v, state: %+v, fsm: %+v", ev, state, fsmData)
	}
	l.Entries = append(l.Entries, LogEntry{
		time.Now(),
		fmt.Sprintf("Rej : %s in %s%s", ev.Name(), state.Name(), detail),
	})
}
