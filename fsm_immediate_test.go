package fsm_test

import (
	"fmt"

	fsm "github.com/johngrange/gofsm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Immediate FSM", func() {
	type fsmData struct {
		followGuardOffToOn bool
	}

	var (
		smb                                    fsm.StateMachineBuilder
		stateMachine                           fsm.ImmediateFSM
		data                                   *fsmData
		onState, offState, startingState, init fsm.StateBuilder
		currStateName                          func() string
		err                                    error
	)

	BeforeEach(func() {
		currStateName = func() string {
			return stateMachine.CurrentState().Name()
		}
		data = &fsmData{}
		smb = fsm.NewFSMBuilder().SetData(data)
		init = smb.GetInitialState()

		startingState = fsm.NewStateBuilder("starting")

		onState = fsm.NewStateBuilder("on")
		offState = fsm.NewStateBuilder("off")

		init.AddTransition(startingState)
		startingState.AddTransition(offState)

		offState.AddTransition(onState).SetEventTrigger("TurnOn")
		onState.AddTransition(offState).SetEventTrigger("TurnOff")
		smb.
			AddState(startingState).
			AddState(onState).
			AddState(offState)

	})
	When("starting an fsm", func() {
		It("should be in the initial state before start is called", func() {
			stateMachine, err = smb.BuildImmediateFSM()
			Expect(err).NotTo(HaveOccurred())
			Expect(stateMachine.CurrentState()).NotTo(BeNil())
			Expect(stateMachine.CurrentState().Name()).To(Equal("initial"))
		})
		It("should transition through unguarded, non event transitions when start is called", func() {
			stateMachine, err = smb.BuildImmediateFSM()
			Expect(err).NotTo(HaveOccurred())
			stateMachine.Start()
			Expect(stateMachine.CurrentState()).NotTo(BeNil())
			Expect(stateMachine.CurrentState().Name()).To(Equal("off"))
		})
	})
	When("dispatching an event", func() {
		It("should not transition if a known event is presented in the wrong state", func() {
			stateMachine, err = smb.BuildImmediateFSM()
			Expect(err).NotTo(HaveOccurred())
			stateMachine.Start()
			Expect(stateMachine.CurrentState().Name()).To(Equal("off"))

			stateMachine.Dispatch(fsm.NewEvent("TurnOff", nil))
			Expect(stateMachine.CurrentState().Name()).To(Equal("off"))
		})
		It("should not transition if an unknown event is presented in the wrong state", func() {
			stateMachine, err = smb.BuildImmediateFSM()
			Expect(err).NotTo(HaveOccurred())
			stateMachine.Start()
			Expect(stateMachine.CurrentState().Name()).To(Equal("off"))

			stateMachine.Dispatch(fsm.NewEvent("nosuch event", nil))
			Expect(stateMachine.CurrentState().Name()).To(Equal("off"))
		})
		It("should not transition if an event is presented in the right state", func() {
			stateMachine, err = smb.BuildImmediateFSM()
			Expect(err).NotTo(HaveOccurred())
			stateMachine.Start()
			Expect(stateMachine.CurrentState().Name()).To(Equal("off"))

			stateMachine.Dispatch(fsm.NewEvent("TurnOn", nil))
			Eventually(currStateName).Should(Equal("on"))
			stateMachine.Dispatch(fsm.NewEvent("TurnOff", nil))
			Eventually(currStateName).Should(Equal("off"))
		})
	})
	When("not started", func() {
		It("should not transition if data changes", func() {
			offState.AddTransition(onState).SetGuard(func(data, eventData interface{}) bool {
				return true
			})
			stateMachine, err = smb.BuildImmediateFSM()
			Expect(err).NotTo(HaveOccurred())
			Expect(stateMachine.CurrentState()).NotTo(BeNil())
			Expect(stateMachine.CurrentState().Name()).To(Equal("initial"))
		})

	})
	When("started", func() {
		It("should not transition if guard evaluates to false", func() {
			offState.AddTransition(onState).SetGuard(func(data, eventData interface{}) bool {
				return false
			})
			stateMachine, err = smb.BuildImmediateFSM()
			Expect(err).NotTo(HaveOccurred())
			stateMachine.Start()
			Consistently(currStateName).Should(Equal("off"))

		})
		It("should not transition if guard evaluates to true", func() {
			offState.AddTransition(onState).SetGuard(func(data, eventData interface{}) bool {
				return true
			})
			stateMachine, err = smb.BuildImmediateFSM()
			Expect(err).NotTo(HaveOccurred())
			stateMachine.Start()
			Eventually(currStateName).Should(Equal("on"))

		})
		It("should not transition when guard fsm data changes", func() {
			// There is no underlying threading, so it should not detect the fsm data changes automatically
			offState.AddTransition(onState).SetGuard(func(data, eventData interface{}) bool {
				d := data.(*fsmData)
				fmt.Fprintf(GinkgoWriter, "%+v", data)
				return d.followGuardOffToOn
			})
			data.followGuardOffToOn = false
			stateMachine, err = smb.BuildImmediateFSM()
			Expect(err).NotTo(HaveOccurred())
			stateMachine.Start()

			Consistently(currStateName).Should(Equal("off"))
			data.followGuardOffToOn = true
			Consistently(currStateName).Should(Equal("off"))

		})
		It("should transition when guard fsm data changes and Tick() is called", func() {
			// There is no underlying threading, so it should not detect the fsm data changes automatically
			offState.AddTransition(onState).SetGuard(func(data, eventData interface{}) bool {
				d := data.(*fsmData)
				fmt.Fprintf(GinkgoWriter, "%+v", data)
				return d.followGuardOffToOn
			})
			data.followGuardOffToOn = false
			stateMachine, err = smb.BuildImmediateFSM()
			Expect(err).NotTo(HaveOccurred())
			stateMachine.Start()

			Consistently(currStateName).Should(Equal("off"))
			data.followGuardOffToOn = true
			Consistently(currStateName).Should(Equal("off"))
			stateMachine.Tick()
			Expect(stateMachine.CurrentState().Name()).To(Equal("on"))
		})

	})
	When("applying visitor pattern", func() {
		It("should visit each element once", func() {
			counter := countingVisitor{}
			stateMachine, err = smb.BuildImmediateFSM()
			Expect(err).NotTo(HaveOccurred())
			stateMachine.Visit(&counter)
			Expect(counter.stateCount).To(Equal(4))
			Expect(counter.transitionCount).To(Equal(4))
			// Check it doesn't change with running the fsm
			stateMachine.Start()
			stateMachine.Dispatch(fsm.NewEvent("TurnOn", nil))
			Expect(stateMachine.CurrentState().Name()).To(Equal("on"))
			stateMachine.Dispatch(fsm.NewEvent("TurnOff", nil))
			Expect(stateMachine.CurrentState().Name()).To(Equal("off"))
			stateMachine.Stop()
			counter2 := countingVisitor{}
			stateMachine.Visit(&counter2)
			Expect(counter).To(Equal(counter2))
		})
	})
})

type countingVisitor struct {
	stateCount      int
	transitionCount int
}

func (c *countingVisitor) VisitState(fsm.State) {
	c.stateCount++
}
func (c *countingVisitor) VisitTransition(fsm.Transition) {
	c.transitionCount++
}
