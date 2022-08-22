package fsm

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FSM Builder", func() {

	When("Building a state", func() {
		It("should return the same object for two build calls", func() {
			stateBuilder := NewStateBuilder("state1")
			s1, err := stateBuilder.build()
			Expect(err).NotTo(HaveOccurred())
			s2, err := stateBuilder.build()
			Expect(err).NotTo(HaveOccurred())
			Expect(s1).To(Equal(s2))
		})
		It("should have the correct number of transitions after buildtransitions", func() {
			sb1 := NewStateBuilder("s1")
			sb2 := NewStateBuilder("s2")
			sb3 := NewStateBuilder("s3")
			sb1.AddTransition(sb2)
			sb2.AddTransition(sb1)
			sb2.AddTransition(sb3)
			s1, err := sb1.build()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(s1.Transitions())).To(BeNumerically("==", 0))
			s2, err := sb2.build()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(s2.Transitions())).To(BeNumerically("==", 0))
			s3, err := sb3.build()
			Expect(len(s3.Transitions())).To(BeNumerically("==", 0))
			Expect(err).NotTo(HaveOccurred())
			err = sb1.buildTransitions()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(s1.Transitions())).To(BeNumerically("==", 1))
			err = sb2.buildTransitions()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(s2.Transitions())).To(BeNumerically("==", 2))
			err = sb1.buildTransitions()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(s3.Transitions())).To(BeNumerically("==", 0))

			// try to build s2 again
			s2_1, err := sb2.build()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(s2_1.Transitions())).To(BeNumerically("==", 2))

		})
	})
	When("building a state machine", func() {
		It("should end up with the correct number of transitions on states", func() {
			smb := NewFSMBuilder()
			sb1 := smb.NewState("s1")
			smb.GetInitialState().AddTransition(sb1)
			sm, err := smb.BuildImmediateFSM()
			Expect(err).NotTo(HaveOccurred())
			smimpl := sm.(*immediateFSMImpl)
			Expect(len(smimpl.states)).To(Equal(2))
			Expect(smimpl.states[0].Name()).To(Equal("initial"))
			Expect(len(smimpl.states[0].Transitions())).To(Equal(1))
		})
	})
})
