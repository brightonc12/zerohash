package vwap

import (
	"log"
	"math/big"
	"sync"
	"zerohash/internal/core/services/numbers"
	"zerohash/internal/domain"
)

// Service - The intent of this Interface and functionality is to distribute Trade on the correct channel that requires
// for it.
type Service interface {
	PublishToSingleCalcChan(calcChan map[string]Calculator, wg *sync.WaitGroup)
}

// impl is the implementation struct of the interface Service
type impl struct {
	tradeChan  chan *domain.Trade
	currencies []string
}

// PublishToSingleCalcChan await for a trade from the tradeChan and delegated the trade to the appropriate channel
func (i *impl) PublishToSingleCalcChan(calcChan map[string]Calculator, wg *sync.WaitGroup) {
	defer func() {
		for _, value := range calcChan {
			close(value.GetTradeChan())
			wg.Done()
		}
	}()

	for t := range i.tradeChan {
		if c, found := calcChan[t.Currency]; found {
			c.GetTradeChan() <- t
		}
	}
}


// New creates a Service type to delegate new trades to the appropriate currency
func New(
	tradeChan chan *domain.Trade,
	currencies []*domain.Currency,
) Service {
	cs := make([]string, len(currencies))

	for _, c := range currencies {
		cs = append(cs, c.Name)
	}
	return &impl{
		tradeChan:  tradeChan,
		currencies: cs,
	}
}

// Calculator is an interface that receives trade and computes the result VWap
type Calculator interface {
	UpdateAgainstTrade(wg *sync.WaitGroup)
	GetTradeChan() chan *domain.Trade
}

// calculator is the implementation of Calculator
// 	tradeChan - chan that contain new updates of the current trade
// 	window - the maximum sliding window to calculate the VWap value against
// 	windowTrades - contains all current trades within the window and it gets updated respectively
// 	curIndex - helps us to track the oldest trade in the windowTrades to add a new one
//	sumPriceQuantity - store the sum of price and quantity multiplied of the windowTrades for easier VWap calculation
//	sumQuantity - stores the sum of quantity of the windowTrades for easier VWap calculation
type calculator struct {
	tradeChan        chan *domain.Trade
	window           int
	windowTrades     []*domain.Trade
	curIndex         int
	sumPriceQuantity *big.Float
	sumQuantity      *big.Float
}


// GetTradeChan - publishes the tradeChan to be consumed out of the calculation struct
func (c *calculator) GetTradeChan() chan *domain.Trade {
	return c.tradeChan
}

// UpdateAgainstTrade - received new trade data from the tradeChan
// updates the sumPriceQuantity, sumQuantity, windowTrades and curIndex
// and prints out the latest VWap value from the given data
func (c *calculator) UpdateAgainstTrade(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	for t := range c.tradeChan {
		newMulPriceQuantity := numbers.NewZeroBigFloat()
		newMulPriceQuantity.Mul(t.Price, t.Quantity)

		c.sumPriceQuantity.
			Add(c.sumPriceQuantity, newMulPriceQuantity)
		c.sumQuantity.
			Add(c.sumQuantity, t.Quantity)

		if  len(c.windowTrades) == c.window && c.windowTrades[c.window - 1] != nil {
			if c.curIndex == c.window {
				c.curIndex = 0
			}
			oldT := c.windowTrades[c.curIndex]

			pMulPriceQuantity := numbers.NewZeroBigFloat()
			pMulPriceQuantity.Mul(oldT.Price, oldT.Quantity)

			c.sumPriceQuantity.
				Sub(c.sumPriceQuantity, pMulPriceQuantity)
			c.sumQuantity.
				Sub(c.sumQuantity, oldT.Quantity)
		}

		result := numbers.NewZeroBigFloat()
		result.Quo(c.sumPriceQuantity, c.sumQuantity)

		c.windowTrades[c.curIndex] = t
		c.curIndex++

		log.Printf("%v vwap: %v", t.Currency, result)
	}
}

// NewCalculator - creates a new Calculator type used to calculate VWap
func NewCalculator(
	tradeChan chan *domain.Trade,
	window int,
) Calculator {

	windowTrades := make([]*domain.Trade, window)

	return &calculator{
		tradeChan:        tradeChan,
		window:           window,
		windowTrades:     windowTrades,
		sumQuantity:      numbers.NewZeroBigFloat(),
		sumPriceQuantity: numbers.NewZeroBigFloat(),
	}
}
