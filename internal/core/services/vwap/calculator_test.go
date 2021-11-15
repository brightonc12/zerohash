package vwap

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"zerohash/internal/core/services/numbers"
	"zerohash/internal/domain"
)

var (
	testWindow = 10
	testTrade  = &domain.Trade{
		Currency: domain.Currencies.EthUsd.Name,
		Price:    numbers.NewBigFloat(40.3),
		Quantity: numbers.NewBigFloat(2),
	}
)

func newTestTrade(currency *domain.Currency) *domain.Trade {
	return &domain.Trade{
		Currency: currency.Name,
		Price:    numbers.NewBigFloat(40.3),
		Quantity: numbers.NewBigFloat(2),
	}
}

func TestImpl_PublishToSingleCalcChan(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{name: "publish to chan success", expected: false},
	}

	currencies := []*domain.Currency{
		domain.Currencies.BtcUsd,
		domain.Currencies.EthBtc,
	}

	testTrades := []*domain.Trade{
		newTestTrade(domain.Currencies.BtcUsd),
		newTestTrade(domain.Currencies.EthBtc),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tradeChan := make(chan *domain.Trade)
			srv := New(tradeChan, currencies)

			calcChan := make(map[string]Calculator)
			for _, c := range currencies {
				newTradeChan := make(chan *domain.Trade)
				calcChan[c.Name] = NewCalculator(newTradeChan, testWindow)
			}
			go func() {
				for _, trd := range testTrades {
					tradeChan <- trd
				}
				close(tradeChan)
			}()

			wg := &sync.WaitGroup{}
			wg.Add(2)
			go srv.PublishToSingleCalcChan(calcChan, wg)

			td := <-calcChan[domain.Currencies.BtcUsd.Name].GetTradeChan()
			assert.Equal(t, td.Currency, domain.Currencies.BtcUsd.Name)

			td = <-calcChan[domain.Currencies.EthBtc.Name].GetTradeChan()
			assert.Equal(t, td.Currency, domain.Currencies.EthBtc.Name)
		})
	}
}

func TestCalculator_GetTradeChan(t *testing.T) {
	tradeChan := make(chan *domain.Trade)
	calc := NewCalculator(tradeChan, 10)

	t.Run("get same private channel ", func(t *testing.T) {
		assert.Equal(t, calc.GetTradeChan(), tradeChan)
	})
}
