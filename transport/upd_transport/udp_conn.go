package updtransport

import (
	"net"
	"time"

	"github.com/andyollylarkin/smudge-custom-transport/transport"
)

type UDPConn struct {
	underlyingConn *net.UDPConn
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (uc *UDPConn) Read(b []byte) (n int, err error) {
	return uc.underlyingConn.Read(b)
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (uc *UDPConn) Write(b []byte) (n int, err error) {
	return uc.underlyingConn.Write(b)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (uc *UDPConn) Close() error {
	return uc.underlyingConn.Close()
}

// LocalAddr returns the local network address, if known.
func (uc *UDPConn) LocalAddr() net.Addr {
	return uc.underlyingConn.LocalAddr()
}

// RemoteAddr returns the remote network address, if known.
func (uc *UDPConn) RemoteAddr() net.Addr {
	return uc.underlyingConn.RemoteAddr()
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
func (uc *UDPConn) SetDeadline(t time.Time) error {
	return uc.underlyingConn.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (uc *UDPConn) SetReadDeadline(t time.Time) error {
	return uc.underlyingConn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (uc *UDPConn) SetWriteDeadline(t time.Time) error {
	return uc.underlyingConn.SetWriteDeadline(t)
}

func (uc *UDPConn) ReadFrom(b []byte) (n int, addr transport.SockAddr, error error) {
	n, udpAddr, err := uc.underlyingConn.ReadFromUDP(b)
	if err != nil {
		return 0, nil, err
	}

	a := &UDPAddr{
		udpAddr,
	}

	return n, a, nil
}
