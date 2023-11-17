package transport

import (
	"net"
)

type SockAddr interface {
	net.Addr
	GetIPAddr() net.IP
	GetPort() int
	GetZone() string
}
