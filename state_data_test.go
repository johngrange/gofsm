package fsm_test

import (
	"fmt"

	fsm "github.com/johngrange/gofsm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tests for data in states", func() {
	type currentCoinPayment struct {
		numCoins  uint
		coinValue uint
	}
	type paymentMeter struct {
		paymentsCollected uint
		ticketsIssued     uint
		ticketCost        uint
		currentPayment    currentCoinPayment
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
		init = fsm.NewState("initial")

		idleState = fsm.NewState("idle")

		acceptingPaymentState = fsm.NewState("acceptingPayment")

		printingTicketState = fsm.NewState("printingTicket")

		init.AddTransition(idleState)

		idleState.AddTransition(acceptingPaymentState).SetTrigger("evInsertCoin").SetAction(
			func(ev fsm.Event, fsmData interface{}) {
				stateData := &(fsmData).(*paymentMeter).currentPayment
				fmt.Fprintf(GinkgoWriter, "stateData: %+v\n", stateData)
				coinAmount := ev.Data().(uint)
				stateData.coinValue += coinAmount
				stateData.numCoins++
				fmt.Fprintf(GinkgoWriter, "stateData: %+v\n", stateData)
			},
		) // parameter is coin value: uint

		acceptingPaymentState.AddTransition(acceptingPaymentState).SetTrigger("evInsertCoin").SetAction(
			func(ev fsm.Event, fsmData interface{}) {
				stateData := &(fsmData).(*paymentMeter).currentPayment
				fmt.Fprintf(GinkgoWriter, "stateData: %+v\n", stateData)

				coinAmount := ev.Data().(uint)
				stateData.coinValue += coinAmount
				stateData.numCoins++

				fmt.Fprintf(GinkgoWriter, "stateData: %+v\n", stateData)
			},
		) // parameter is coin value: uint

		acceptingPaymentState.AddTransition(printingTicketState).SetTrigger("evPrintTicket").
			SetGuard(func(fsmData, eventData interface{}) bool {
				meterData := (fsmData).(*paymentMeter)

				return meterData.currentPayment.coinValue >= meterData.ticketCost
			}).
			SetAction(func(ev fsm.Event, fsmData interface{}) {
				meter := (fsmData).(*paymentMeter)
				fmt.Fprintf(GinkgoWriter, "Printing ticket for %dp\n", meter.currentPayment.coinValue)
				meter.paymentsCollected += meter.currentPayment.coinValue
				meter.ticketsIssued++
			})

		acceptingPaymentState.OnExit(func(state fsm.FSMState, fsmData interface{}) {
			fmt.Fprintf(GinkgoWriter, "onExit")
			meterData := (fsmData).(*paymentMeter)
			meterData.currentPayment.coinValue = 0
			meterData.currentPayment.numCoins = 0
		})

		printingTicketState.AddTransition(idleState)

		paymentMeterSM = fsm.NewImmediateFSM(init, paymentMeterData)
		paymentMeterSM.
			AddState(idleState).
			AddState(acceptingPaymentState).
			AddState(printingTicketState)

	})
	When("putting several coins in to equal ticket value", func() {

		It("should allow a ticket to be issued, and the meter data should reflect that", func() {
			const coinAmount = 50
			const customers = 10
			const coinsPerCustomer = 6
			paymentMeterSM.Start()

			for customer := 0; customer < customers; customer++ {
				Expect(paymentMeterSM.CurrentState().Name()).To(Equal("idle"))
				for coin := 0; coin < coinsPerCustomer; coin++ {
					fmt.Fprintf(GinkgoWriter, "inserting coin %d\n", coin)
					paymentMeterSM.Dispatch(fsm.NewEvent("evInsertCoin", uint(coinAmount)))
					Expect(paymentMeterSM.CurrentState().Name()).To(Equal("acceptingPayment"))
				}
				Expect(paymentMeterSM.CurrentState().Name()).To(Equal("acceptingPayment"))
				paymentMeterSM.Dispatch(fsm.NewEvent("evPrintTicket", nil))
				Expect(paymentMeterSM.CurrentState().Name()).To(Equal("idle"))
			}
			meterData := paymentMeterSM.GetData().(*paymentMeter)
			Expect(meterData.ticketsIssued).To(BeNumerically("==", customers))
			Expect(meterData.paymentsCollected).To(BeNumerically("==", coinAmount*coinsPerCustomer*customers))
		})
	})

})
