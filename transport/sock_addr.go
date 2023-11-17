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

type UDPAddr struct {
	*net.UDPAddr
}

func (ua *UDPAddr) IPAddr() net.IP {
	return ua.IP
}

func (ua *UDPAddr) GetPort() int {
	return ua.Port
}

func (ua *UDPAddr) GetZone() string {
	return ua.Zone
}
