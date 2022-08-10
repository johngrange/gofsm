# README
Very simple finite state machine in golang.

Based, loosely on FSM from Rational Rhapsody.

Usage: see gofsm_threaded_test.go

## Goals and Non Goals
Goals:
- To provide useful, simple to use FSM for golang
- To follow UML state diagram semantics as far as reasonable for the features implemented.
- Provide concurrent and non-concurrent implementations.
Non Goals:
- Full implementation of UML state machine specification