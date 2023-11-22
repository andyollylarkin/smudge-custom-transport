package internal

import (
	"fmt"
	"net"
	"net/http"
	"strings"

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
	realIpArr := r.Header.Get("X-Forwarded-For")
	realIps := strings.Split(realIpArr, ",")
	// if X-Forwarded-For header is set, we behind nginx. Set connection remote address to X-Forwarded-For
	if len(realIps) > 1 { // if realIps == 1, then client not set self ip in header. real ip is equal nginx $remote_addr
		_, port, err := net.SplitHostPort(wsconn.RemoteAddr().String())
		if err != nil {
			return nil, fmt.Errorf("cant assign real remote addr: %w", err)
		}

		addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(realIpArr, port))
		if err != nil {
			return nil, fmt.Errorf("cant assign real remote addr: %w", err)
		}

		adapter := NewWsConn(wsconn, addr)

		return adapter, nil
	}

	return nil, fmt.Errorf("request doesnt contains X-Forwarded-For header")
}

func createConnNoRequest(wsconn *websocket.Conn) (*WsConnAdapter, error) {
	return NewWsConn(wsconn, wsconn.RemoteAddr()), nil
}
