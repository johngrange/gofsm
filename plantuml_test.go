package fsm_test

import (
	"bytes"
	"fmt"
	"os"
	"path"

	fsm "github.com/johngrange/gofsm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testOutputDir = "./test-output"

var _ = Describe("Plant UML Rendering", func() {
	type fsmData struct {
	}

	var (
		stateMachine                           fsm.FSM
		data                                   *fsmData
		onState, offState, startingState, init fsm.StateBuilder
		err                                    error
	)

	When("rendering uml", func() {
		It("should get the output right!", func() {
			data = &fsmData{}
			stateMachineBuilder := fsm.NewFSMBuilder().SetData(data)
			init = stateMachineBuilder.GetInitialState()

			startingState = fsm.NewStateBuilder("starting")
			startingState.OnEntry(func(state fsm.State, fsmData interface{}, dispatcher fsm.Dispatcher) {}, "initialise system")

			onState = fsm.NewStateBuilder("on", "the on state")
			offState = fsm.NewStateBuilder("off")

			init.AddTransition(startingState)
			startingState.AddTransition(offState)

			offState.AddTransition(onState).SetEventTrigger("TurnOn").SetGuard(func(fsmData, eventData interface{}) bool { return true }, "power==active")
			onState.AddTransition(offState).SetEventTrigger("TurnOff").SetEffect(func(ev fsm.Event, fsmData interface{}, dispatcher fsm.Dispatcher) {}, "perform effect")

			onState.AddTransition(stateMachineBuilder.AddFinalState()).SetEventTrigger("FatalError").SetEffect(func(ev fsm.Event, fsmData interface{}, dispatcher fsm.Dispatcher) {}, "panic!")

			onState.OnExit(func(state fsm.State, fsmData interface{}, dispatcher fsm.Dispatcher) {}, "turn out lights")
			fmt.Fprintf(GinkgoWriter, "fsm: %+v, %T\n", stateMachine, stateMachine)

			stateMachineBuilder.
				AddState(startingState).
				AddState(onState).
				AddState(offState)

			stateMachine, err = stateMachineBuilder.BuildThreadedFSM()
			Expect(err).NotTo(HaveOccurred())
			fmt.Fprintf(GinkgoWriter, "fsm: %+v, %T\n", stateMachine, stateMachine)

			buf := bytes.Buffer{}

			err = fsm.RenderPlantUML(&buf, stateMachine)
			Expect(err).NotTo(HaveOccurred())

			expectedUML :=
				`@startuml
[*] --> starting
starting : entry/initialise system
starting --> off
on : the on state
on : exit/turn out lights
on --> off : TurnOff/perform effect
on --> [*] : FatalError/panic!
off --> on : TurnOn [power==active] 
@enduml
`
			fmt.Fprintf(GinkgoWriter, "%s\n", buf.String())
			Expect(buf.String()).To(Equal(expectedUML))
			err = os.MkdirAll(testOutputDir, 0755)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(path.Join(testOutputDir, "testone.uml"), buf.Bytes(), 0600)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
