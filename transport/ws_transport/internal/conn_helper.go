package internal

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
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
	realIp, ok := r.Header["X-Real-IP"]
	// if X-Real-IP header is set, we behind nginx. Set connection remote address to X-Real-IP
	if ok {
		_, port, err := net.SplitHostPort(wsconn.RemoteAddr().String())
		if err != nil {
			return nil, fmt.Errorf("cant assign real remote addr: %w", err)
		}

		addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(realIp[0], port))
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
