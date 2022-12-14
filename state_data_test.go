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
		idleState, acceptingPaymentState, printingTicketState, init fsm.StateBuilder
		err                                                         error
	)

	BeforeEach(func() {

		// car park payment meter model

		paymentMeterData := &paymentMeter{
			ticketCost: 300,
		}

		smb := fsm.NewFSMBuilder().SetData(paymentMeterData)

		init = smb.GetInitialState()

		idleState = fsm.NewStateBuilder("idle")

		acceptingPaymentState = fsm.NewStateBuilder("acceptingPayment")

		printingTicketState = fsm.NewStateBuilder("printingTicket")

		init.AddTransition(idleState)

		idleState.AddTransition(acceptingPaymentState).SetEventTrigger("evInsertCoin").SetEffect(
			func(ev fsm.Event, fsmData interface{}, dispatcher fsm.Dispatcher) {
				stateData := &(fsmData).(*paymentMeter).currentPayment
				fmt.Fprintf(GinkgoWriter, "stateData: %+v\n", stateData)
				coinAmount := ev.Data().(uint)
				stateData.coinValue += coinAmount
				stateData.numCoins++
				fmt.Fprintf(GinkgoWriter, "stateData: %+v\n", stateData)
			},
		) // parameter is coin value: uint

		acceptingPaymentState.AddTransition(acceptingPaymentState).SetEventTrigger("evInsertCoin").SetEffect(
			func(ev fsm.Event, fsmData interface{}, dispatcher fsm.Dispatcher) {
				stateData := &(fsmData).(*paymentMeter).currentPayment
				fmt.Fprintf(GinkgoWriter, "stateData: %+v\n", stateData)

				coinAmount := ev.Data().(uint)
				stateData.coinValue += coinAmount
				stateData.numCoins++

				fmt.Fprintf(GinkgoWriter, "stateData: %+v\n", stateData)
			},
		) // parameter is coin value: uint

		acceptingPaymentState.AddTransition(printingTicketState).SetEventTrigger("evPrintTicket").
			SetGuard(func(fsmData, eventData interface{}) bool {
				meterData := (fsmData).(*paymentMeter)

				return meterData.currentPayment.coinValue >= meterData.ticketCost
			}).
			SetEffect(func(ev fsm.Event, fsmData interface{}, dispatcher fsm.Dispatcher) {
				meter := (fsmData).(*paymentMeter)
				fmt.Fprintf(GinkgoWriter, "Printing ticket for %dp\n", meter.currentPayment.coinValue)
				meter.paymentsCollected += meter.currentPayment.coinValue
				meter.ticketsIssued++
			})

		acceptingPaymentState.OnExit(func(state fsm.State, fsmData interface{}, dispatcher fsm.Dispatcher) {
			fmt.Fprintf(GinkgoWriter, "onExit\n")
			meterData := (fsmData).(*paymentMeter)
			meterData.currentPayment.coinValue = 0
			meterData.currentPayment.numCoins = 0
		})

		printingTicketState.AddTransition(idleState)

		smb.
			AddState(idleState).
			AddState(acceptingPaymentState).
			AddState(printingTicketState)
		paymentMeterSM, err = smb.BuildImmediateFSM()
		Expect(err).NotTo(HaveOccurred())

	})
	When("putting several coins in to equal ticket value", func() {

		It("should allow a ticket to be issued, and the meter data should reflect that", func() {
			const coinAmount = 50
			const customers = 10
			const coinsPerCustomer = 6
			paymentMeterSM.Start()

			for customer := 0; customer < customers; customer++ {
				fmt.Fprintf(GinkgoWriter, "serving customer %d\n", customer)
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
