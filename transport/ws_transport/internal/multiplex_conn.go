package internal

import (
	"net"
	"sync"
	"time"

	"github.com/andyollylarkin/smudge-custom-transport"
	"github.com/andyollylarkin/smudge-custom-transport/transport"
)

type readData struct {
	readed     int
	readedFrom transport.SockAddr
	err        error
	data       []byte
}

type MultiplexConn struct {
	laddr         transport.SockAddr
	wg            sync.WaitGroup
	dataChan      chan readData
	onCloseChan   chan struct{}
	connChan      chan transport.GenericConn
	connErrorChan chan net.Addr
	logger        smudge.Logger
}

func NewMuxConn(laddr transport.SockAddr, logger smudge.Logger) (*MultiplexConn, chan net.Addr) {
	connErrChan := make(chan net.Addr)
	c := &MultiplexConn{
		laddr:         laddr,
		wg:            sync.WaitGroup{},
		dataChan:      make(chan readData),
		onCloseChan:   make(chan struct{}),
		connChan:      make(chan transport.GenericConn),
		connErrorChan: connErrChan,
		logger:        logger,
	}

	go c.handleLoop()

	return c, connErrChan
}

func (mc *MultiplexConn) HandleNewConn(conn transport.GenericConn) {
	mc.connChan <- conn
}

func (mc *MultiplexConn) handleLoop() {
	for {
		select {
		case conn := <-mc.connChan:
			go mc.handleRead(conn)
		case <-mc.onCloseChan:
			return
		}
	}
}

func (mc *MultiplexConn) handleRead(conn transport.GenericConn) {
	buf := make([]byte, 11)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			mc.connErrorChan <- conn.RemoteAddr()

			return
		}

		tcpaddr, err := net.ResolveTCPAddr("tcp", conn.RemoteAddr().String())

		mc.dataChan <- readData{
			readed:     n,
			readedFrom: &WsAddr{WsAddrTCP: *tcpaddr},
			err:        err,
			data:       buf,
		}

		buf = make([]byte, 11)
	}
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (mc *MultiplexConn) Read(b []byte) (n int, err error) {
	panic("not implemented") // TODO: Implement
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (mc *MultiplexConn) Write(b []byte) (n int, err error) {
	panic("not implemented") // TODO: Implement
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (mc *MultiplexConn) Close() error {
	close(mc.onCloseChan)

	return nil
}

// LocalAddr returns the local network address, if known.
func (mc *MultiplexConn) LocalAddr() net.Addr {
	return mc.laddr
}

// RemoteAddr returns the remote network address, if known.
func (mc *MultiplexConn) RemoteAddr() net.Addr {
	return nil
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
func (mc *MultiplexConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (mc *MultiplexConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (mc *MultiplexConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (mc *MultiplexConn) ReadFrom(b []byte) (n int, addr transport.SockAddr, error error) {
	data := <-mc.dataChan

	copy(b, data.data)

	mc.logger.Logf(smudge.LogDebug, "Read %v from %s", data.data, data.readedFrom.GetIPAddr())

	return data.readed, data.readedFrom, data.err
}
