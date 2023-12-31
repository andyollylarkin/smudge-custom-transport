package internal

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/andyollylarkin/smudge-custom-transport/transport"
	"github.com/gorilla/websocket"
)

// WsConnAdapter adapter (wrapper) around web socket connection for correspond net.Conn interface.
type WsConnAdapter struct {
	SocketConn     *websocket.Conn
	realRemoteAddr net.Addr
	mu             sync.Mutex
}

func NewWsConn(conn *websocket.Conn, realRemoteAddr net.Addr) *WsConnAdapter {
	wsc := new(WsConnAdapter)
	wsc.SocketConn = conn
	wsc.realRemoteAddr = realRemoteAddr

	return wsc
}

func (wsc *WsConnAdapter) ReadFrom(b []byte) (n int, addr transport.SockAddr, err error) {
	n, err = wsc.Read(b)
	if err != nil {
		return
	}

	sockAddr := wsc.SocketConn.RemoteAddr()

	tcpAddr, ok := sockAddr.(*net.TCPAddr)

	if !ok {
		return 0, nil, fmt.Errorf("cant resolve ws tcp addr. RemoteAddr() socket is not TCP Addr")
	}

	wsAddr := &WsAddr{WsAddrTCP: *tcpAddr}

	return n, wsAddr, nil
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (wsc *WsConnAdapter) Read(b []byte) (n int, err error) {
	t, r, err := wsc.SocketConn.NextReader()
	if err != nil {
		return 0, err
	}

	if t != websocket.TextMessage {
		return 0, fmt.Errorf("invalid websocket message type. Require text message")
	}

	return r.Read(b)
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (wsc *WsConnAdapter) Write(b []byte) (n int, err error) {
	wsc.mu.Lock()
	defer wsc.mu.Unlock()

	err = wsc.SocketConn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return 0, err
	}

	return len(b), err
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
// **NOTICE** this function do nothing, because we don't want the connection to be closed after every data transfer.
func (wsc *WsConnAdapter) Close() error {
	return nil
}

// ActuallyClose actyally close websocket connection.
func (wsc *WsConnAdapter) ActuallyClose() error {
	return wsc.SocketConn.Close()
}

// LocalAddr returns the local network address, if known.
func (wsc *WsConnAdapter) LocalAddr() net.Addr {
	return wsc.SocketConn.LocalAddr()
}

// RemoteAddr returns the remote network address, if known.
func (wsc *WsConnAdapter) RemoteAddr() net.Addr {
	return wsc.realRemoteAddr
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail instead of blocking. The deadline applies to all future
// and pending I/O, not just the immediately following call to
// Read or Write. After a deadline has been exceeded, the
// connection can be refreshed by setting a deadline in the future.
//
// If the deadline is exceeded a call to Read or Write or to other
// I/O methods will return an error that wraps os.ErrDeadlineExceeded.
// This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
// The error's Timeout method will return true, but note that there
// are other possible errors for which the Timeout method will
// return true even if the deadline has not been exceeded.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (wsc *WsConnAdapter) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (wsc *WsConnAdapter) SetReadDeadline(t time.Time) error {
	return wsc.SocketConn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (wsc *WsConnAdapter) SetWriteDeadline(t time.Time) error {
	return wsc.SocketConn.SetWriteDeadline(t)
}
