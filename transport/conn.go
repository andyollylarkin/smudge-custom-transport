package transport

import "net"

type GenericConn interface {
	net.Conn
	ReadFrom(b []byte) (n int, addr SockAddr, error error)
}
