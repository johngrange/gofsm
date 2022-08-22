package fsm_test

import (
	"fmt"

	fsm "github.com/johngrange/gofsm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Threaded FSM", func() {
	type fsmData struct {
		errorDuringOn      bool
		followGuardOffToOn bool
	}

	var (
		stateMachine                           fsm.FSM
		data                                   *fsmData
		onState, offState, startingState, init fsm.StateBuilder
		currStateName                          func() string
		err                                    error
		stateMachineBuilder                    fsm.FSMBuilder
	)

	BeforeEach(func() {
		currStateName = func() string {
			return stateMachine.CurrentState().Name()
		}
		data = &fsmData{}
		stateMachineBuilder = fsm.NewFSMBuilder()
		stateMachineBuilder.SetData(data)
		init = stateMachineBuilder.GetInitialState()

		startingState = fsm.NewStateBuilder("starting")

		onState = fsm.NewStateBuilder("on")
		offState = fsm.NewStateBuilder("off")

		init.AddTransition(startingState)
		startingState.AddTransition(offState)

		offState.AddTransition(onState).SetTrigger("TurnOn")
		onState.AddTransition(offState).SetTrigger("TurnOff")
		fmt.Fprintf(GinkgoWriter, "fsm: %+v, %T\n", stateMachine, stateMachine)

		stateMachineBuilder.
			AddState(startingState).
			AddState(onState).
			AddState(offState)

	})
	When("starting an fsm", func() {
		It("should be in the initial state before start is called", func() {
			stateMachine, err = stateMachineBuilder.BuildThreadedFSM()
			Expect(err).NotTo(HaveOccurred())
			fmt.Fprintf(GinkgoWriter, "fsm: %+v, %T\n", stateMachine, stateMachine)
			Expect(stateMachine.CurrentState()).NotTo(BeNil())
			Expect(stateMachine.CurrentState().Name()).To(Equal("initial"))
		})
		It("should transition through unguarded, non event transitions when start is called", func() {
			stateMachine, err = stateMachineBuilder.BuildThreadedFSM()
			Expect(err).NotTo(HaveOccurred())
			fmt.Fprintf(GinkgoWriter, "fsm: %+v, %T\n", stateMachine, stateMachine)
			stateMachine.Start()
			Expect(stateMachine.CurrentState()).NotTo(BeNil())
			Expect(stateMachine.CurrentState().Name()).To(Equal("off"))
		})
	})
	When("dispatching an event", func() {
		It("should not transition if a known event is presented in the wrong state", func() {
			stateMachine, err = stateMachineBuilder.BuildThreadedFSM()
			Expect(err).NotTo(HaveOccurred())
			fmt.Fprintf(GinkgoWriter, "fsm: %+v, %T\n", stateMachine, stateMachine)
			stateMachine.Start()
			Expect(stateMachine.CurrentState().Name()).To(Equal("off"))

			stateMachine.Dispatch(fsm.NewEvent("TurnOff", nil))
			Expect(stateMachine.CurrentState().Name()).To(Equal("off"))
		})
		It("should not transition if an unknown event is presented in the wrong state", func() {
			stateMachine, err = stateMachineBuilder.BuildThreadedFSM()
			Expect(err).NotTo(HaveOccurred())
			fmt.Fprintf(GinkgoWriter, "fsm: %+v, %T\n", stateMachine, stateMachine)
			stateMachine.Start()
			Expect(stateMachine.CurrentState().Name()).To(Equal("off"))

			stateMachine.Dispatch(fsm.NewEvent("nosuch event", nil))
			Expect(stateMachine.CurrentState().Name()).To(Equal("off"))
		})
		It("should not transition if an event is presented in the right state", func() {
			stateMachine, err = stateMachineBuilder.BuildThreadedFSM()
			Expect(err).NotTo(HaveOccurred())
			fmt.Fprintf(GinkgoWriter, "fsm: %+v, %T\n", stateMachine, stateMachine)
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
			fmt.Fprintf(GinkgoWriter, "fsm: %+v, %T\n", stateMachine, stateMachine)
			offState.AddTransition(onState).SetGuard(func(data, eventData interface{}) bool {
				return true
			})
			stateMachine, err = stateMachineBuilder.BuildThreadedFSM()
			Expect(err).NotTo(HaveOccurred())
			Expect(stateMachine.CurrentState()).NotTo(BeNil())
			Expect(stateMachine.CurrentState().Name()).To(Equal("initial"))
		})

	})
	When("started", func() {
		It("should not transition if guard evaluates to false", func() {
			fmt.Fprintf(GinkgoWriter, "fsm: %+v, %T\n", stateMachine, stateMachine)
			offState.AddTransition(onState).SetGuard(func(data, eventData interface{}) bool {
				return false
			})
			stateMachine, err = stateMachineBuilder.BuildThreadedFSM()
			Expect(err).NotTo(HaveOccurred())
			stateMachine.Start()
			Consistently(currStateName).Should(Equal("off"))

		})
		It("should transition if guard evaluates to true", func() {
			fmt.Fprintf(GinkgoWriter, "fsm: %+v, %T\n", stateMachine, stateMachine)
			offState.AddTransition(onState).SetGuard(func(data, eventData interface{}) bool {
				return true
			})
			stateMachine, err = stateMachineBuilder.BuildThreadedFSM()
			Expect(err).NotTo(HaveOccurred())
			stateMachine.Start()
			Eventually(currStateName).Should(Equal("on"))

		})
		It("should transition when guard test changes", func() {

			offState.AddTransition(onState).SetGuard(func(data, eventData interface{}) bool {
				d := data.(*fsmData)
				fmt.Fprintf(GinkgoWriter, "%+v", data)
				return d.followGuardOffToOn
			})
			stateMachine, err = stateMachineBuilder.BuildThreadedFSM()
			Expect(err).NotTo(HaveOccurred())
			data.followGuardOffToOn = false
			stateMachine.Start()

			Consistently(currStateName).Should(Equal("off"))
			data.followGuardOffToOn = true
			Eventually(currStateName).Should(Equal("on"))

		})

	})
})
