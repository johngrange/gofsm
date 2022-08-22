package fsm_test

import (
	"fmt"
	"time"

	fsm "github.com/johngrange/gofsm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("test tracers", func() {
	type fsmData struct {
	}

	var (
		stateMachine                           fsm.ImmediateFSM
		data                                   *fsmData
		onState, offState, startingState, init fsm.StateBuilder
		err                                    error
	)

	BeforeEach(func() {
		data = &fsmData{}
		smb := fsm.NewFSMBuilder().SetData(data)
		init = smb.GetInitialState()

		startingState = fsm.NewStateBuilder("starting")

		onState = fsm.NewStateBuilder("on")
		offState = fsm.NewStateBuilder("off")

		init.AddTransition(startingState)
		startingState.AddTransition(offState)

		offState.AddTransition(onState).SetTrigger("TurnOn")
		onState.AddTransition(offState).SetTrigger("TurnOff")
		smb.
			AddState(startingState).
			AddState(onState).
			AddState(offState)
		stateMachine, err = smb.BuildImmediateFSM()
		Expect(err).NotTo(HaveOccurred())
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

			Expect(counter.StateCounts).To(Equal(map[string]uint64{
				"initial":  1,
				"starting": 1,
				"off":      2,
				"on":       1,
			}))
			Expect(counter.RejectedEventCounts).To(Equal(map[string]uint64{
				"TurnOff": 1,
				"TurnOn":  1,
			}))
		})
	})
	When("using a log tracer", func() {
		It("should log states correctly", func() {
			logger := fsm.NewFSMLogger()
			defer func() {
				for _, l := range logger.Entries {
					fmt.Fprintln(GinkgoWriter, l)
				}
			}()

			stateMachine.AddTracer(logger)
			startTime := time.Now()
			stateMachine.Start()                                // Should go (En)initial(Ex)->(T)(En)starting(Ex)->(T)(En)off: 7
			stateMachine.Dispatch(fsm.NewEvent("TurnOff", nil)) // should not transition 1
			stateMachine.Dispatch(fsm.NewEvent("TurnOn", nil))  // off(Ex)->(T)(En)on: 3
			stateMachine.Dispatch(fsm.NewEvent("TurnOn", nil))  // should not transition 1
			stateMachine.Dispatch(fsm.NewEvent("TurnOff", nil)) // on(Ex)->(T)(En)off: 3

			endTime := time.Now()
			Expect(len(logger.Entries)).To(BeNumerically("==", 15))

			for _, l := range logger.Entries {
				Expect(l.When).To(BeTemporally(">=", startTime))
				Expect(l.When).To(BeTemporally("<=", endTime))
				Expect(len(l.Message)).To(BeNumerically(">", 0))
			}
		})
	})
})
