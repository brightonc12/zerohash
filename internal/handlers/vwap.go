package vwapHdl

import (
	"context"
	"log"
	"sync"
	"zerohash/internal/core/services/coinbase"
	"zerohash/internal/core/services/vwap"
	"zerohash/internal/domain"
)

var (
	window = 200
)

func RunVWapAgainstTrade(ctx context.Context, wg *sync.WaitGroup) {

	currencies := []*domain.Currency{domain.Currencies.BtcUsd, domain.Currencies.EthUsd, domain.Currencies.EthBtc}
	tr := coinbase.NewTraderReader(currencies)

	tradeChan := make(chan *domain.Trade)

	err := tr.Connect()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	wg.Add(1)
	go tr.ReadTradeToChan(ctx, tradeChan, wg)

	if err := tr.Subscribe(); err != nil {
		log.Fatalf("error subscribing to trade reader : %v", err)
	}

	calc := vwap.New(tradeChan, currencies)

	calcChan := make(map[string]vwap.Calculator)
	for _, c := range currencies {
		newTradeChan := make(chan *domain.Trade)
		calcChan[c.Name] = vwap.NewCalculator(newTradeChan, window)
	}

	wg.Add(1)
	go calc.PublishToSingleCalcChan(calcChan, wg)
	////
	for _, value := range calcChan {
		wg.Add(1)
		go value.UpdateAgainstTrade(wg)
	}

	//defer func() {
	//	log.Println("closes")
	//	err := tr.Unsubscribe()
	//	if err != nil {
	//		log.Fatalf("Error: %v", err)
	//	}
	//}()

}
