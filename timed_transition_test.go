package fsm_test

import (
	"fmt"
	"time"

	fsm "github.com/johngrange/gofsm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Timed transition tests", func() {
	printLog := func(logger *fsm.Logger) {
		fmt.Fprintf(GinkgoWriter, "Printing log\n")
		if len(logger.Entries) == 0 {
			return
		}
		firstTime := logger.Entries[0].When
		for _, entry := range logger.Entries {
			fmt.Fprintf(GinkgoWriter, "%s: %s\n", entry.When.Sub(firstTime), entry.Message)
		}
	}
	Describe("simple state machine", func() {
		type fsmData struct {
			followGuardOnToOff bool
		}
		var (
			data                    *fsmData
			onState, offState, init fsm.StateBuilder
			stateMachineBuilder     fsm.StateMachineBuilder
			logger                  *fsm.Logger
		)

		BeforeEach(func() {

			data = &fsmData{}
			stateMachineBuilder = fsm.NewFSMBuilder()
			stateMachineBuilder.SetData(data)
			init = stateMachineBuilder.GetInitialState()

			onState = fsm.NewStateBuilder("on")
			offState = fsm.NewStateBuilder("off")

			init.AddTransition(offState)

			offState.AddTransition(onState).SetTimedTrigger(time.Millisecond * 100)
			onState.AddTransition(offState).SetTimedTrigger(time.Millisecond * 150).SetGuard(func(data, eventData interface{}) bool {
				d := data.(*fsmData)
				fmt.Fprintf(GinkgoWriter, "checking on-off guard with data %+v\n", d)
				return d.followGuardOnToOff
			})

			stateMachineBuilder.
				AddState(onState).
				AddState(offState)
			logger = fsm.NewFSMLogger()

		})
		AfterEach(func() {
			printLog(logger)
		})
		When("processing timers on threaded fsm", func() {
			It("should transition without external stimuli", func() {
				stateMachine, err := stateMachineBuilder.BuildThreadedFSM()
				Expect(err).NotTo(HaveOccurred())

				stateMachine.AddTracer(logger)

				currStateName := func() string {
					return stateMachine.CurrentState().Name()
				}
				stateMachine.Start()
				Eventually(currStateName, "200ms").Should(Equal("off"))
				Eventually(currStateName, "500ms").Should(Equal("on"))
				Consistently(currStateName, "100ms").Should(Equal("on"))
				data.followGuardOnToOff = true
				Eventually(currStateName, "200ms").Should(Equal("off"))
			})
		})
		When("processing timers on immediate fsm", func() {
			It("should only transition during Tick() calls", func() {
				stateMachine, err := stateMachineBuilder.BuildImmediateFSM()
				Expect(err).NotTo(HaveOccurred())

				stateMachine.AddTracer(logger)

				currStateName := func() string {
					return stateMachine.CurrentState().Name()
				}
				stateMachine.Start()
				Expect(currStateName()).To(Equal("off"))
				Consistently(currStateName, "200ms").Should(Equal("off"))
				stateMachine.Tick()
				Expect(currStateName()).To(Equal("on"))
				Consistently(currStateName, "200ms").Should(Equal("on"))
				stateMachine.Tick()
				Consistently(currStateName, "200ms").Should(Equal("on"))
				data.followGuardOnToOff = true
				Consistently(currStateName, "200ms").Should(Equal("on"))
				stateMachine.Tick()
				Expect(currStateName()).To(Equal("off"))
				Consistently(currStateName, "200ms").Should(Equal("off"))
				stateMachine.Tick()
				Expect(currStateName()).To(Equal("on"))
			})
		})
	})
	Describe("testing for multiple timers expiring together", func() {
		type fsmData struct {
			abGuard bool
			acGuard bool
		}
		var (
			data                   *fsmData
			stateA, stateB, stateC fsm.StateBuilder
			stateMachineBuilder    fsm.StateMachineBuilder
			logger                 *fsm.Logger
		)

		BeforeEach(func() {

			data = &fsmData{}
			stateMachineBuilder = fsm.NewFSMBuilder()
			stateMachineBuilder.SetData(data)

			stateA = fsm.NewStateBuilder("stateA")
			stateB = fsm.NewStateBuilder("stateB")
			stateC = fsm.NewStateBuilder("stateC")

			stateMachineBuilder.GetInitialState().AddTransition(stateA)

			stateA.AddTransition(stateB).SetTimedTrigger(time.Millisecond * 200).SetGuard(func(data, eventData interface{}) bool {
				return data.(*fsmData).abGuard
			})
			stateA.AddTransition(stateC).SetTimedTrigger(time.Millisecond * 100).SetGuard(func(data, eventData interface{}) bool {
				return data.(*fsmData).acGuard
			})

			stateMachineBuilder.
				AddState(stateA).
				AddState(stateB).
				AddState(stateC)
			logger = fsm.NewFSMLogger()
		})
		AfterEach(func() {
			printLog(logger)
		})
		When("processing timers on threaded fsm", func() {
			It("should follow a longer timer if shorter guard is false", func() {
				stateMachine, err := stateMachineBuilder.BuildThreadedFSM()
				Expect(err).NotTo(HaveOccurred())

				stateMachine.AddTracer(logger)

				currStateName := func() string {
					return stateMachine.CurrentState().Name()
				}
				data.abGuard = true
				stateMachine.Start()
				Eventually(currStateName, "1ms").Should(Equal("stateA"))
				Eventually(currStateName, "500ms").Should(Equal("stateB"))
			})
			It("should follow a shorter timer if both guards are true", func() {
				stateMachine, err := stateMachineBuilder.BuildThreadedFSM()
				Expect(err).NotTo(HaveOccurred())

				stateMachine.AddTracer(logger)

				currStateName := func() string {
					return stateMachine.CurrentState().Name()
				}
				data.abGuard = true
				data.acGuard = true
				stateMachine.Start()
				Eventually(currStateName, "1ms").Should(Equal("stateA"))
				Eventually(currStateName, "300ms").Should(Equal("stateC"))
				// Give time for the ab timer to have occurred
				Consistently(currStateName, "500ms").Should(Equal("stateC"))
			})
		})
		When("processing timers on immediate fsm", func() {
			It("should follow a longer timer if shorter guard is false", func() {
				stateMachine, err := stateMachineBuilder.BuildImmediateFSM()
				Expect(err).NotTo(HaveOccurred())

				stateMachine.AddTracer(logger)

				data.abGuard = true
				data.acGuard = true
				stateMachine.Start()
				Expect(stateMachine.CurrentState().Name()).To(Equal("stateA"))
				time.Sleep(100 * time.Millisecond)
				stateMachine.Tick()
				Expect(stateMachine.CurrentState().Name()).To(Equal("stateC"))
				time.Sleep(200 * time.Millisecond)
				stateMachine.Tick()
				Expect(stateMachine.CurrentState().Name()).To(Equal("stateC"))
			})
		})

	})
})
