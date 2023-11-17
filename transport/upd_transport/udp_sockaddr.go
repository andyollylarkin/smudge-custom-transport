package updtransport

import "net"

type UDPAddr struct {
	*net.UDPAddr
}

func (ua *UDPAddr) GetIPAddr() net.IP {
	return ua.IP
}

func (ua *UDPAddr) GetPort() int {
	return ua.Port
}

func (ua *UDPAddr) GetZone() string {
	return ua.Zone
}
