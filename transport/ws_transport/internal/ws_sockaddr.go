package internal

import "net"

type WsAddr struct {
	WsAddrTCP net.TCPAddr
}

func (ws *WsAddr) Network() string {
	return ws.WsAddrTCP.Network()
}

func (ws *WsAddr) String() string {
	return ws.WsAddrTCP.String()
}

func (ws *WsAddr) GetIPAddr() net.IP {
	return ws.WsAddrTCP.IP
}

func (ws *WsAddr) GetPort() int {
	return ws.WsAddrTCP.Port
}

func (ws *WsAddr) GetZone() string {
	return ws.WsAddrTCP.Zone
}
