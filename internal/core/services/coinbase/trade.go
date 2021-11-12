package coinbase

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"sync"
	"zerohash/internal/core/services/numbers"
	"zerohash/internal/domain"
)

var (
	socketUrl = "wss://ws-feed.exchange.coinbase.com"



	)

type TraderReader interface {
	Connect() error
	Subscribe() error
	ReadTradeToChan(
		ctx context.Context,
		tradeChan chan *domain.Trade,
		wg *sync.WaitGroup,
	)
	Unsubscribe() error
}

type reader struct {
	conn       *websocket.Conn
	currencies []*domain.Currency
}

func (r *reader) Connect() error {
	ws, err := websocket.Dial(socketUrl, "", "http://localhost")
	if err != nil {
		return fmt.Errorf("error while connecting to coinbase websocket: %w", err)
	}
	r.conn = ws
	return nil
}

func (r *reader) Subscribe() error {

	var pIds []string

	for _, c := range r.currencies {
		pIds = append(pIds, c.Name)
	}

	subMsg := &domain.CoinbaseSubscribeMsg{
		Type: "subscribe",
		Channels:    []string{"matches"},
		ProductIDs: pIds,
	}

	subEvt, _ := json.Marshal(subMsg)

	if _, err := r.conn.Write(subEvt); err != nil {
		return fmt.Errorf("error subscribing to coinbase websocket: %w", err)
	}

	return nil
}

func (r *reader) ReadTradeToChan(
	ctx context.Context,
	tradeChan chan *domain.Trade,
	wg *sync.WaitGroup,
) {
	if r.conn == nil {
		log.Fatalf("coinbase websocket is not yet connected")
	}

	defer func() {
		log.Println("finished ReadTradeToChan")
		close(tradeChan)
		wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return

		default:
			m := &domain.MatchMsg{}

			if err := json.NewDecoder(r.conn).Decode(&m); err != nil {
				log.Printf("error while decoding message: %v", err)
				continue
			}

			if m.TradeID == 0 {
				log.Printf("trade container no values %+v", m)
				continue
			}

			t, err := r.extractTradeFromMsg(m)

			if err != nil {
				log.Printf("Error extracting match message to trade")
				continue
			}
			tradeChan <- t
		}
	}
}

func (r *reader) extractTradeFromMsg(match *domain.MatchMsg) (*domain.Trade, error) {
	p, err := numbers.ParseBigFloat(match.Price)
	if err != nil {
		return nil, err
	}

	s, err := numbers.ParseBigFloat(match.Size)
	if err != nil {
		return nil, err
	}

	t := &domain.Trade{
		Currency: match.ProductID,
		Price:    p,
		Quantity: s,
	}

	return t, nil
}

func (r *reader) Unsubscribe() error {
	unsubMsg := &domain.CoinbaseUnsubscribeMsg{
		Type: "unsubscribe",
		Channels: []string{"matches"},
	}
	unsubEvt, _ := json.Marshal(unsubMsg)

	if _, err := r.conn.Write(unsubEvt); err != nil {
		return fmt.Errorf("error unsubscribing to coinbase websocket: %w", err)
	}
	return nil
}

func NewTraderReader(currencies []*domain.Currency) TraderReader {
	return &reader{
		currencies: currencies,
	}
}
