package fsm_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	fsm "github.com/johngrange/gofsm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testOutputDir = "./test-output"

var _ = Describe("Plant UML Rendering", func() {
	type fsmData struct {
		errorDuringOn      bool
		followGuardOffToOn bool
	}

	var (
		stateMachine                           fsm.FSMBuilder
		data                                   *fsmData
		onState, offState, startingState, init fsm.StateBuilder
	)

	When("rendering uml", func() {
		It("should get the output right!", func() {
			data = &fsmData{}
			init = fsm.NewState("initial")

			startingState = fsm.NewState("starting")
			startingState.OnEntry(func(state fsm.State, fsmData interface{}, dispatcher fsm.Dispatcher) {}, "initialise system")

			onState = fsm.NewState("on", "the on state")
			offState = fsm.NewState("off")

			init.AddTransition(startingState)
			startingState.AddTransition(offState)

			offState.AddTransition(onState).SetTrigger("TurnOn").SetGuard(func(fsmData, eventData interface{}) bool { return true }, "power==active")
			onState.AddTransition(offState).SetTrigger("TurnOff").SetEffect(func(ev fsm.Event, fsmData interface{}, dispatcher fsm.Dispatcher) {}, "perform effect")
			onState.OnExit(func(state fsm.State, fsmData interface{}, dispatcher fsm.Dispatcher) {}, "turn out lights")
			stateMachine = fsm.NewThreadedFSM(init, data)
			fmt.Fprintf(GinkgoWriter, "fsm: %+v, %T\n", stateMachine, stateMachine)

			stateMachine.
				AddState(startingState).
				AddState(onState).
				AddState(offState)

			fmt.Fprintf(GinkgoWriter, "fsm: %+v, %T\n", stateMachine, stateMachine)

			buf := bytes.Buffer{}

			err := fsm.RenderPlantUML(&buf, stateMachine)
			Expect(err).NotTo(HaveOccurred())

			expectedUML :=
				`@startuml
[*] --> initial
initial --> starting
starting : entry/initialise system
starting --> off
on : the on state
on : exit/turn out lights
on --> off : TurnOff/perform effect
off --> on : TurnOn [power==active] 
@enduml
`
			fmt.Fprintf(GinkgoWriter, "%s\n", string(buf.Bytes()))
			Expect(string(buf.Bytes())).To(Equal(expectedUML))
			os.MkdirAll(testOutputDir, 0755)
			err = ioutil.WriteFile(path.Join(testOutputDir, "testone.uml"), buf.Bytes(), 0644)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
