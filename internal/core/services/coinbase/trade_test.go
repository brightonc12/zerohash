package coinbase

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"zerohash/internal/domain"
)

func TestNewTraderReader(t *testing.T) {
	t.Run("should return valid Trade Reader", func(t *testing.T) {
		r := NewTraderReader([]*domain.Currency{})
		if _, ok := r.(TraderReader); !ok {
			t.Errorf("NewTraderReader must be of TraderReader type")
		}
	})
}

func TestReader_Connect(t *testing.T) {
	tests := []struct {
		name      string
		returnErr bool
	}{
		{"connection failed", true},
		{"connection succeed", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setUpWSServerTest(handleServerRequest(t, tt.returnErr))
			defer s.Close()
			r := NewTraderReader([]*domain.Currency{})

			if err := r.Connect(); err != nil {
				if tt.returnErr {
					return
				}
				t.Errorf("test failed")
			}
		})
	}
}

func TestReader_Subscribe(t *testing.T) {
	tests := []struct {
		name      string
		returnErr bool
	}{
		{"subscribe fail", true},
		{"subscribe succeed", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currencies := []*domain.Currency{
				domain.Currencies.EthBtc,
				domain.Currencies.BtcUsd,
				domain.Currencies.EthUsd,
			}
			done := make(chan bool)

			s := setUpWSServerTest(WSRequest(t, func(conn *websocket.Conn) {
				if !tt.returnErr {
					m := &domain.CoinbaseSubscribeMsg{}
					err := conn.ReadJSON(m)
					assert.NoError(t, err)
					assert.Equal(t, "subscribe", m.Type)
					assert.Equal(t, []string{"matches"}, m.Channels)
					assert.Equal(t, len(currencies), len(m.ProductIDs))
				}
			}, done))

			defer s.Close()
			r := NewTraderReader(currencies)
			r.Connect()

			if tt.returnErr {
				r.Close()
			} else {
				defer r.Close()
			}

			err := r.Subscribe()
			if tt.returnErr {
				assert.Error(t, err)
				<-done
			} else {
				<-done
				assert.NoError(t, err)
			}
		})
	}
}

func TestReader_ReadTradeToChan(t *testing.T) {

	tests := []struct {
		name string
	}{
		{"s"},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			s := setUpWSServerTest(successRequest(t))
			defer s.Close()
			r := NewTraderReader([]*domain.Currency{})
			wg := &sync.WaitGroup{}
			tradeChan := make(chan *domain.Trade)
			ctx := context.Background()

			r.ReadTradeToChan(ctx, tradeChan, wg)
		})
	}
}

func TestReader_Unsubscribe(t *testing.T) {
	tests := []struct {
		name      string
		returnErr bool
	}{
		{"unsubscribe fail", true},
		{"unsubscribe succeed", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currencies := []*domain.Currency{
				domain.Currencies.EthBtc,
				domain.Currencies.BtcUsd,
				domain.Currencies.EthUsd,
			}
			done := make(chan bool)

			s := setUpWSServerTest(WSRequest(t, func(conn *websocket.Conn) {
				if !tt.returnErr {
					m := &domain.CoinbaseUnsubscribeMsg{}
					err := conn.ReadJSON(m)
					assert.NoError(t, err)
					assert.Equal(t, "unsubscribe", m.Type)
					assert.Equal(t, []string{"matches"}, m.Channels)
				}
			}, done))

			defer s.Close()
			r := NewTraderReader(currencies)
			r.Connect()

			if tt.returnErr {
				r.Close()
			} else {
				defer r.Close()
			}

			err := r.Unsubscribe()
			if tt.returnErr {
				assert.Error(t, err)
				<-done
			} else {
				<-done
				assert.NoError(t, err)
			}
		})
	}
}

func setUpWSServerTest(handlerFunc http.HandlerFunc) *httptest.Server {
	s := httptest.NewServer(handlerFunc)
	socketUrl = "ws" + strings.TrimPrefix(s.URL, "http")
	return s
}

func handleServerRequest(t *testing.T, fail bool) func(w http.ResponseWriter, r *http.Request) {
	if fail {
		return failRequest()
	}
	return successRequest(t)
}

func failRequest() func(w http.ResponseWriter, r *http.Request) {
	// will cause an error
	return func(w http.ResponseWriter, r *http.Request) {}
}

func successRequest(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := u.Upgrade(w, r, nil)
		assert.NoError(t, err)
		defer conn.Close()
	}
}

func WSRequest(
	t *testing.T,
	callback func(conn *websocket.Conn),
	done chan bool,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := u.Upgrade(w, r, nil)
		assert.NoError(t, err)
		defer conn.Close()

		callback(conn)
		done <- true
	}
}
