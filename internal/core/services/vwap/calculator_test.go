package vwap

import (
	"testing"
	"zerohash/internal/core/services/numbers"
	"zerohash/internal/domain"
)

var (
	testWindow = 10
	testTrade = &domain.Trade{
		Currency: domain.Currencies.EthUsd.Name,
		Price: numbers.NewBigFloat(40.3),
		Quantity: numbers.NewBigFloat(2),
	}
)

func TestImpl_PublishToSingleCalcChan(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		from     string
		to       string
		expected bool
	}{
		{name: "before time", input: "Sunday 1AM - 1PM", from: "9AM", to: "5PM", expected: false},
	}

	tradeChan := make(chan *domain.Trade)
	currencies := []*domain.Currency{
		domain.Currencies.EthBtc,
		domain.Currencies.BtcUsd,
	}

	_ = New(tradeChan, currencies)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			calcChan := make(map[string]Calculator)
			for _, c := range currencies {
				newTradeChan := make(chan *domain.Trade)
				calcChan[c.Name] = NewCalculator(newTradeChan, testWindow)
			}
			//srv.PublishToSingleCalcChan()
		})

	}
}

func TestCalculator_UpdateAgainstTrade(t *testing.T) {
}

func TestCalculator_GetTradeChan(t *testing.T) {
}