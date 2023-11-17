package internal

import (
	"net"
	"net/netip"
)

type UDPAddr struct {
	Uaddr *net.UDPAddr
}

func (ua *UDPAddr) GetIPAddr() net.IP {
	return ua.Uaddr.IP
}

func (ua *UDPAddr) GetPort() int {
	return ua.Uaddr.Port
}

func (ua *UDPAddr) GetZone() string {
	return ua.Uaddr.Zone
}

func (ua *UDPAddr) AddrPort() netip.AddrPort {
	return ua.Uaddr.AddrPort()
}

func (ua *UDPAddr) Network() string {
	return ua.Uaddr.Network()
}

func (ua *UDPAddr) String() string {
	return ua.Uaddr.String()
}
