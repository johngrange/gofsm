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
		stateMachine                           fsm.FSM
		data                                   *fsmData
		onState, offState, startingState, init fsm.FSMState
	)

	When("rendering uml", func() {
		It("should get the output right!", func() {
			data = &fsmData{}
			init = fsm.NewState("initial")

			startingState = fsm.NewState("starting")

			onState = fsm.NewState("on")
			offState = fsm.NewState("off")

			init.AddTransition(startingState)
			startingState.AddTransition(offState)

			offState.AddTransition(onState).SetTrigger("TurnOn")
			onState.AddTransition(offState).SetTrigger("TurnOff")
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
starting --> off
on --> off : TurnOff
off --> on : TurnOn
@enduml
`
			Expect(string(buf.Bytes())).To(Equal(expectedUML))
			os.MkdirAll(testOutputDir, 0755)
			err = ioutil.WriteFile(path.Join(testOutputDir, "testone.uml"), buf.Bytes(), 0644)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
