package transport

import "context"

type Transport interface {
	Listen(network string, addr SockAddr) (GenericConn, error)
	Dial(ctx context.Context, laddr SockAddr, raddr SockAddr) (GenericConn, error)
	ResolveAddr(network, addr string) (SockAddr, error)
	AllowMulticast() bool
	// Return network, udp, websockets, tcp, ipv4, etc.
	Network() string
}
