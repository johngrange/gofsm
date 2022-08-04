package fsm_test

import (
	"time"

	fsm "github.com/johngrange/gofsm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("test tracers", func() {
	type fsmData struct {
		followGuardOffToOn bool
	}

	var (
		stateMachine                           fsm.FSM
		data                                   *fsmData
		onState, offState, startingState, init fsm.FSMState
	)

	BeforeEach(func() {
		data = &fsmData{}
		init = fsm.NewState("initial")

		startingState = fsm.NewState("starting")

		onState = fsm.NewState("on")
		offState = fsm.NewState("off")

		init.AddTransition(startingState)
		startingState.AddTransition(offState)

		offState.AddTransition(onState).SetTrigger("TurnOn")
		onState.AddTransition(offState).SetTrigger("TurnOff")
		stateMachine = fsm.NewImmediateFSM(init, data).
			AddState(startingState).
			AddState(onState).
			AddState(offState)

	})
	When("using a count tracer", func() {
		It("should count states correctly", func() {
			counter := fsm.NewStateCounter()
			stateMachine.AddTracer(counter)
			stateMachine.Start()                                // Should go initial->starting->off
			stateMachine.Dispatch(fsm.NewEvent("TurnOff", nil)) // should not transition
			stateMachine.Dispatch(fsm.NewEvent("TurnOn", nil))  // off->on
			stateMachine.Dispatch(fsm.NewEvent("TurnOn", nil))  // should not transition
			stateMachine.Dispatch(fsm.NewEvent("TurnOff", nil)) // on->off

			Expect(counter.Counts["initial"]).To(BeNumerically("==", 1))
			Expect(counter.Counts["starting"]).To(BeNumerically("==", 1))
			Expect(counter.Counts["off"]).To(BeNumerically("==", 2))
			Expect(counter.Counts["on"]).To(BeNumerically("==", 1))
		})
	})
	When("using a log tracer", func() {
		It("should log states correctly", func() {
			logger := fsm.NewFSMLogger()
			stateMachine.AddTracer(logger)
			startTime := time.Now()
			stateMachine.Start()                                // Should go (En)initial(Ex)->(T)(En)starting(Ex)->(T)(En)off: 7
			stateMachine.Dispatch(fsm.NewEvent("TurnOff", nil)) // should not transition
			stateMachine.Dispatch(fsm.NewEvent("TurnOn", nil))  // off(Ex)->(T)(En)on: 3
			stateMachine.Dispatch(fsm.NewEvent("TurnOn", nil))  // should not transition
			stateMachine.Dispatch(fsm.NewEvent("TurnOff", nil)) // on(Ex)->(T)(En)off: 3
			endTime := time.Now()
			Expect(len(logger.Entries)).To(BeNumerically("==", 13))

			for _, l := range logger.Entries {
				Expect(l.When).To(BeTemporally(">=", startTime))
				Expect(l.When).To(BeTemporally("<=", endTime))
				Expect(len(l.Message)).To(BeNumerically(">", 0))
			}
		})
	})
})
