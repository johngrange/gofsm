package fsm_test

import (
	"fmt"
	"reflect"

	fsm "github.com/johngrange/gofsm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tests for data in states", func() {
	type paymentMeter struct {
		paymentsCollected uint
		ticketsIssued     uint
		ticketCost        uint
	}

	var (
		paymentMeterSM                                              fsm.ImmediateFSM
		idleState, acceptingPaymentState, printingTicketState, init fsm.FSMState
	)

	BeforeEach(func() {

		// car park payment meter model

		paymentMeterData := &paymentMeter{
			ticketCost: 300,
		}
		type currentCoinPayment struct {
			numCoins  uint
			coinValue uint
		}
		init = fsm.NewState("initial")

		idleState = fsm.NewState("idle")

		acceptingPaymentState = fsm.NewState("acceptingPayment")
		acceptingPaymentState.SetDataFactory(func() interface{} {
			fmt.Fprint(GinkgoWriter, "data factory called\n")
			return &currentCoinPayment{}
		})

		printingTicketState = fsm.NewState("printingTicket")

		init.AddTransition(idleState)

		idleState.AddTransition(acceptingPaymentState).SetTrigger("evInsertCoin").SetAction(
			func(fromState, toState fsm.FSMState, ev fsm.Event, fsmData interface{}) {
				fmt.Fprintf(GinkgoWriter, "csd: %+v, %v\n", toState.GetCurrentData(), reflect.TypeOf(toState.GetCurrentData()))
				stateData := (toState.GetCurrentData()).(*currentCoinPayment)
				coinAmount := ev.Data().(uint)
				stateData.coinValue += coinAmount
				stateData.numCoins++
				fmt.Fprintf(GinkgoWriter, "stateData: %+v\n", stateData)
			},
		) // parameter is coin value: uint

		acceptingPaymentState.AddTransition(acceptingPaymentState).SetTrigger("evInsertCoin").SetAction(
			func(fromState, toState fsm.FSMState, ev fsm.Event, fsmData interface{}) {
				fmt.Fprintf(GinkgoWriter, "from, to: %s, -> %s\n", fromState.Name(), toState.Name())
				fmt.Fprintf(GinkgoWriter, "csd: %+v, %v\n", toState.GetCurrentData(), reflect.TypeOf(toState.GetCurrentData()))
				stateData := (toState.GetCurrentData()).(*currentCoinPayment)
				coinAmount := ev.Data().(uint)
				stateData.coinValue += coinAmount
				stateData.numCoins++

				fmt.Fprintf(GinkgoWriter, "stateData: %+v\n", stateData)
			},
		) // parameter is coin value: uint

		acceptingPaymentState.AddTransition(printingTicketState).SetTrigger("evPrintTicket").
			SetGuard(func(state fsm.FSMState, fsmData, eventData interface{}) bool {
				fmt.Fprintf(GinkgoWriter, "current state in guard: %s\n", state.Name())
				currentStateData := (state.GetCurrentData()).(*currentCoinPayment)
				meterData := (fsmData).(*paymentMeter)

				return currentStateData.coinValue >= meterData.ticketCost
			}).
			SetAction(func(fromState, toState fsm.FSMState, ev fsm.Event, fsmData interface{}) {
				meter := (fsmData).(*paymentMeter)
				collectionStateData := (fromState.GetCurrentData()).(*currentCoinPayment)
				fmt.Fprintf(GinkgoWriter, "Printing ticket for %dp\n", collectionStateData.coinValue)
				meter.paymentsCollected += collectionStateData.coinValue
				meter.ticketsIssued++
			})

		printingTicketState.AddTransition(idleState)

		paymentMeterSM = fsm.NewImmediateFSM(init, paymentMeterData)
		paymentMeterSM.
			AddState(idleState).
			AddState(acceptingPaymentState).
			AddState(printingTicketState)

	})
	FWhen("putting several coins in to equal ticket value", func() {

		It("should allow a ticket to be issued, and the meter data should reflect that", func() {
			const coinAmount = 50
			const customers = 10
			const coinsPerCustomer = 3
			paymentMeterSM.Start()

			// for customer := 0; customer < customers; customer++ {
			Expect(paymentMeterSM.CurrentState().Name()).To(Equal("idle"))
			for coin := 0; coin < coinsPerCustomer; coin++ {
				fmt.Fprintf(GinkgoWriter, "inserting coin %d\n", coin)
				paymentMeterSM.Dispatch(fsm.NewEvent("evInsertCoin", uint(coinAmount)))
				Expect(paymentMeterSM.CurrentState().Name()).To(Equal("acceptingPayment"))
			}
			Expect(paymentMeterSM.CurrentState().Name()).To(Equal("acceptingPayment"))
			// paymentMeterSM.Dispatch(fsm.NewEvent("evPrintTicket", nil))
			Expect(paymentMeterSM.CurrentState().Name()).To(Equal("idle"))
			// }
			meterData := paymentMeterSM.GetData().(*paymentMeter)
			Expect(meterData.ticketsIssued).To(Equal(customers))
			Expect(meterData.paymentsCollected).To(Equal(coinAmount * coinsPerCustomer * customers))
		})
	})

})
