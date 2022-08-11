package fsm_test

import (
	"fmt"

	fsm "github.com/johngrange/gofsm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tests for dispatching events inside the state machine", func() {
	When("Using an immediate FSM", func() {
		It("should work when an event dispatched in trigger", func() {
			ctr := fsm.NewStateCounter()
			on := fsm.NewState("on")
			off := fsm.NewState("off")
			errstate := fsm.NewState("error")

			off.AddTransition(on).SetTrigger("on").SetEffect(func(ev fsm.Event, fsmData interface{}, dispatcher fsm.Dispatcher) {
				dispatcher.Dispatch(fsm.NewEvent("error", nil))
			})

			on.AddTransition(off).SetTrigger("off")

			on.AddTransition(errstate).SetTrigger("error")

			sm := fsm.NewImmediateFSM(off, nil)
			sm.AddTracer(ctr)
			sm.Start()
			Expect(sm.CurrentState().Name()).To(Equal("off"))
			sm.Dispatch(fsm.NewEvent("on", nil))
			Expect(sm.CurrentState().Name()).To(Equal("error"))
			Expect(ctr.StateCounts).To(Equal(map[string]uint64{
				"off":   1,
				"on":    1,
				"error": 1,
			}))

		})
		It("should work when an event dispatched in entry and exit", func() {
			ctr := fsm.NewStateCounter()
			on := fsm.NewState("on")
			off := fsm.NewState("off")
			errstate := fsm.NewState("error")
			fixing := fsm.NewState("fixing")
			off.AddTransition(errstate).SetTrigger("error")

			off.AddTransition(on).SetTrigger("on")

			on.AddTransition(off).SetTrigger("off")
			on.OnExit(func(state fsm.State, fsmData interface{}, dispatcher fsm.Dispatcher) {
				dispatcher.Dispatch(fsm.NewEvent("error", nil))
			})

			errstate.OnEntry(func(state fsm.State, fsmData interface{}, dispatcher fsm.Dispatcher) {
				dispatcher.Dispatch(fsm.NewEvent("fixit", nil))
			})

			errstate.AddTransition(fixing).SetTrigger("fixit")
			sm := fsm.NewImmediateFSM(off, nil)
			sm.AddTracer(ctr)
			sm.Start()
			Expect(sm.CurrentState().Name()).To(Equal("off"))
			sm.Dispatch(fsm.NewEvent("on", nil))
			Expect(sm.CurrentState().Name()).To(Equal("on"))
			sm.Dispatch(fsm.NewEvent("off", nil))
			Expect(sm.CurrentState().Name()).To(Equal("fixing"))
			Expect(ctr.StateCounts).To(Equal(map[string]uint64{
				"off":    2,
				"on":     1,
				"error":  1,
				"fixing": 1,
			}))
		})
	})
	When("Using an threaded FSM", func() {
		It("should work when an event dispatched in trigger", func() {
			ctr := fsm.NewStateCounter()
			on := fsm.NewState("on")
			off := fsm.NewState("off")
			errstate := fsm.NewState("error")

			off.AddTransition(on).SetTrigger("on").SetEffect(func(ev fsm.Event, fsmData interface{}, dispatcher fsm.Dispatcher) {
				dispatcher.Dispatch(fsm.NewEvent("error", nil))
			})

			on.AddTransition(off).SetTrigger("off")

			on.AddTransition(errstate).SetTrigger("error")

			sm := fsm.NewThreadedFSM(off, nil)
			sm.AddTracer(ctr)
			sm.Start()
			defer sm.Stop()
			Expect(sm.CurrentState().Name()).To(Equal("off"))
			sm.Dispatch(fsm.NewEvent("on", nil))
			Eventually(func() string { return sm.CurrentState().Name() }).Should(Equal("error"))
			Expect(ctr.StateCounts).To(Equal(map[string]uint64{
				"off":   1,
				"on":    1,
				"error": 1,
			}))
			Expect(ctr.RejectedEventCounts).To(Equal(map[string]uint64{}))

		})
		It("should work when an event dispatched in entry and exit", func() {
			lg := fsm.NewFSMLogger()
			ctr := fsm.NewStateCounter()
			on := fsm.NewState("on")
			off := fsm.NewState("off")
			errstate := fsm.NewState("error")
			fixing := fsm.NewState("fixing")
			off.AddTransition(errstate).SetTrigger("error")

			off.AddTransition(on).SetTrigger("on")

			on.AddTransition(off).SetTrigger("off")
			on.OnExit(func(state fsm.State, fsmData interface{}, dispatcher fsm.Dispatcher) {
				dispatcher.Dispatch(fsm.NewEvent("error", nil))
			})

			errstate.OnEntry(func(state fsm.State, fsmData interface{}, dispatcher fsm.Dispatcher) {
				dispatcher.Dispatch(fsm.NewEvent("fixit", nil))
			})

			errstate.AddTransition(fixing).SetTrigger("fixit")
			sm := fsm.NewThreadedFSM(off, nil)
			sm.AddTracer(ctr)
			sm.AddTracer(lg)
			defer func() {
				for _, l := range lg.Entries {
					fmt.Fprintln(GinkgoWriter, l)
				}
			}()
			sm.Start()
			defer sm.Stop()
			Expect(sm.CurrentState().Name()).To(Equal("off"))
			sm.Dispatch(fsm.NewEvent("on", nil))
			Eventually(func() string { return sm.CurrentState().Name() }).Should(Equal("on"))
			sm.Dispatch(fsm.NewEvent("off", nil))
			Eventually(func() string { return sm.CurrentState().Name() }).Should(Equal("fixing"))
			Expect(ctr.StateCounts).To(Equal(map[string]uint64{
				"off":    2,
				"on":     1,
				"error":  1,
				"fixing": 1,
			}))
			Expect(ctr.RejectedEventCounts).To(Equal(map[string]uint64{}))
		})
	})
})
