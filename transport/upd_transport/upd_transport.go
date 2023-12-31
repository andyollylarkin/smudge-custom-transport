package updtransport

import (
	"context"
	"net"

	"github.com/andyollylarkin/smudge-custom-transport/transport"
	"github.com/andyollylarkin/smudge-custom-transport/transport/upd_transport/internal"
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

	genericConn := &internal.UDPConn{
		UnderlyingConn: udpConn,
	}

	return genericConn, nil
}

func (ut *UDPTransport) Dial(ctx context.Context, laddr transport.SockAddr,
	raddr transport.SockAddr,
) (transport.GenericConn, error) {
	var ludpAddr *net.UDPAddr

	var rudpAddr *net.UDPAddr

	var err error

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

	genericConn := &internal.UDPConn{
		UnderlyingConn: udpConn,
	}

	return genericConn, nil
}

func (ut *UDPTransport) ResolveAddr(network string, addr string) (transport.SockAddr, error) {
	sockAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	sa := &internal.UDPAddr{
		Uaddr: sockAddr,
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
