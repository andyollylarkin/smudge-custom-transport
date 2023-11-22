package internal

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	minPortRange = 10_000
	maxPortRange = 65_535
)

// NewWsConnAdapter create new websocket conn adapter. If req is nil -> return remote socket addr as is.
// If requst not nil and contains X-Real-IP, return adapter with replaced remote address.
func NewWsConnAdapter(r *http.Request, wsconn *websocket.Conn) (*WsConnAdapter, error) {
	if r != nil {
		return createConnWithRequest(r, wsconn)
	}

	return createConnNoRequest(wsconn)
}

func createConnWithRequest(r *http.Request, wsconn *websocket.Conn) (*WsConnAdapter, error) {
	realIp := r.Header.Get("X-Real-IP")
	// if X-Real-IP header is set, we behind nginx. Set connection remote address to X-Real-IP
	if realIp != "" {
		// generate random port for unique detect hosts with same ip behind same gateway
		port := generateRandomPort()

		addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(realIp, port))
		if err != nil {
			return nil, fmt.Errorf("cant assign real remote addr: %w", err)
		}

		adapter := NewWsConn(wsconn, addr)

		return adapter, nil
	}

	return nil, fmt.Errorf("request doesnt contains X-Real-IP header")
}

func createConnNoRequest(wsconn *websocket.Conn) (*WsConnAdapter, error) {
	return NewWsConn(wsconn, wsconn.RemoteAddr()), nil
}

type IntRange struct {
	min, max int
}

// get next random value within the interval including min and max
func (ir *IntRange) NextRandom(r *rand.Rand) int {
	return r.Intn(ir.max-ir.min+1) + ir.min
}

func generateRandomPort() string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ir := IntRange{minPortRange, maxPortRange}

	return fmt.Sprintf("%d", ir.NextRandom(r))
}
