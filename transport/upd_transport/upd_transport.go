package updtransport

import (
	"context"
	"net"

	"github.com/andyollylarkin/smudge-custom-transport/transport"
)

type UDPTransport struct{}

func (ut *UDPTransport) Listen(network string, addr transport.SockAddr) (transport.GenericConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr.String())
	if err != nil {
		return nil, err
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	genericConn := &UDPConn{
		underlyingConn: udpConn,
	}

	return genericConn, nil
}

func (ut *UDPTransport) Dial(ctx context.Context, laddr transport.SockAddr,
	raddr transport.SockAddr,
) (transport.GenericConn, error) {
	var ludpAddr *net.UDPAddr

	var rudpAddr *net.UDPAddr

	var err error

	if laddr != nil {
		ludpAddr, err = net.ResolveUDPAddr("udp", laddr.String())
		if err != nil {
			return nil, err
		}
	}

	if raddr != nil {
		rudpAddr, err = net.ResolveUDPAddr("udp", raddr.String())
		if err != nil {
			return nil, err
		}
	}

	udpConn, err := net.DialUDP("udp", ludpAddr, rudpAddr)
	if err != nil {
		return nil, err
	}

	genericConn := &UDPConn{
		underlyingConn: udpConn,
	}

	return genericConn, nil
}

func (ut *UDPTransport) ResolveAddr(network string, addr string) (transport.SockAddr, error) {
	sockAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	sa := &UDPAddr{
		sockAddr,
	}

	return sa, nil
}

func (ut *UDPTransport) AllowMulticast() bool {
	return true
}

// Return network, udp, websockets, tcp, ipv4, etc.
func (ut *UDPTransport) Network() string {
	return "udp"
}
