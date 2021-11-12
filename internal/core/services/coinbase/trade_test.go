package coinbase

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
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
			done := make(chan bool, 1)

			s := setUpWSServerTest(WSRequest(t, func(conn *websocket.Conn) {
				if !tt.returnErr {
					m := &domain.CoinbaseSubscribeMsg{}
					err := conn.ReadJSON(m)
					assert.NoError(t, err)
					assert.Equal(t, "subscribe", m.Type)
					assert.Equal(t, []string{"matches"}, m.Channels)
					assert.Equal(t, 3, len(m.ProductIDs))
				}
				done <- true
			}))

			defer s.Close()
			r := NewTraderReader([]*domain.Currency{})

			err := r.Subscribe()
			if tt.returnErr {
				assert.Error(t, err)
				<- done
			} else {
				<- done
				assert.NoError(t, err)
			}
		})
	}
}

func TestReader_ReadTradeToChan(t *testing.T) {
}

func TestReader_Unsubscribe(t *testing.T) {

}

func TestExtractTradeFromMsg(t *testing.T) {

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

func WSRequest(t *testing.T, callback func(conn *websocket.Conn)) func(w http.ResponseWriter, r *http.Request) {
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
	}
}
