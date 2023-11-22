package wstransport

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/andyollylarkin/smudge-custom-transport"
	"github.com/andyollylarkin/smudge-custom-transport/transport"
	"github.com/andyollylarkin/smudge-custom-transport/transport/ws_transport/internal"
	"github.com/gorilla/websocket"
)

const (
	MaxLRUCacheItems int = 100
)

var (
	upgrader websocket.Upgrader
)

type WsTransport struct {
	cache              *internal.ConnectionStore
	wg                 sync.WaitGroup
	listenIp           net.IP
	remoteWsServerPort *int
	wsBasePath         string
	connChan           chan *internal.WsConnAdapter
	logger             smudge.Logger
}

func NewWsTransport(logger smudge.Logger, remoteWsServerPort *int, wsBasePath string) (*WsTransport, error) {
	cache := internal.NewConnectionStore()

	t := new(WsTransport)

	t.logger = logger
	t.remoteWsServerPort = remoteWsServerPort
	t.wsBasePath = wsBasePath
	t.connChan = make(chan *internal.WsConnAdapter)

	t.cache = cache

	return t, nil
}

// UpgageWebsocket upgrade http request to websocket connection. Pass it to web server handler.
func (wst *WsTransport) UpgageWebsocket(w http.ResponseWriter, r *http.Request) error {
	wsconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("cant upgrade websocket connection: %w", err)
	}

	_, ok, err := wst.cache.ConnCacheGet(wsconn.RemoteAddr())
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	adapter, err := internal.NewWsConnAdapter(r, wsconn)
	if err != nil {
		return err
	}

	err = wst.cache.ConnCacheSet(adapter.RemoteAddr(), adapter)
	if err != nil {
		return err
	}

	wst.connChan <- adapter

	return nil
}

func (wst *WsTransport) Listen(network string, addr transport.SockAddr) (transport.GenericConn, error) {
	wst.logger.Log(smudge.LogWarn,
		"using websocket transport. Some features not working properly (multicast, message sending.)")

	muxConn, connErrChan := internal.NewMuxConn(addr, wst.logger)

	go func() {
		for c := range wst.connChan {
			muxConn.HandleNewConn(c)
		}
	}()

	go wst.connCloseMonitor(connErrChan)

	// TODO: listen on close and then close all opened connections

	return muxConn, nil
}

func (wst *WsTransport) connCloseMonitor(connErrChan chan net.Addr) {
	for addr := range connErrChan {
		conn, ok, err := wst.cache.ConnCacheGet(addr)
		if err != nil || !ok {
			continue
		}

		conn.ActuallyClose()

		wst.cache.ConnCacheRemove(addr)

		wst.logger.Logf(smudge.LogDebug, "Actually close %s", conn.RemoteAddr().String())
	}
}

func (wst *WsTransport) Dial(ctx context.Context, laddr transport.SockAddr,
	raddr transport.SockAddr,
) (transport.GenericConn, error) {
	c, ok, err := wst.cache.ConnCacheGet(raddr)
	if err != nil {
		return nil, err
	}

	// return cached connection
	if ok {
		return c, nil
	}

	if raddr == nil {
		return nil, fmt.Errorf("invalid addr format for raddr. should be host:port, or host. Given nil")
	}

	var remoteAddr string

	if wst.remoteWsServerPort != nil {
		ip, _, err := net.SplitHostPort(raddr.String())
		if err != nil {
			return nil, err
		}

		remoteAddr = net.JoinHostPort(ip, strconv.Itoa(*wst.remoteWsServerPort))
	} else {
		remoteAddr = raddr.String()
	}

	var basePath string

	if wst.wsBasePath == "" {
		basePath = WebsocketRoutePath
	} else {
		basePath = wst.wsBasePath
	}

	url := url.URL{
		Scheme: "ws",
		Host:   remoteAddr,
		Path:   basePath,
	}

	header, err := wst.getDialHeaders()
	if err != nil {
		return nil, err
	}

	wsconn, _, err := websocket.DefaultDialer.Dial(url.String(), header)
	if err != nil {
		return nil, err
	}

	adapter, err := internal.NewWsConnAdapter(nil, wsconn)
	if err != nil {
		return nil, err
	}

	err = wst.cache.ConnCacheSet(adapter.RemoteAddr(), adapter)
	if err != nil {
		return nil, err
	}

	wst.connChan <- adapter

	return adapter, nil
}

func (wst *WsTransport) getDialHeaders() (http.Header, error) {
	var err error
	if wst.listenIp == nil {
		wst.listenIp, err = smudge.TryGetLocalIPv4()
		if err != nil {
			return http.Header{}, err
		}
	}

	header := http.Header{}
	h, _, err := net.SplitHostPort(wst.listenIp.String())
	if err != nil {
		return header, err
	}
	header.Set("X-Forwarded-For", h)

	return header, err
}

func (wst *WsTransport) ResolveAddr(network string, addr string) (transport.SockAddr, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	wsa := &internal.WsAddr{
		WsAddrTCP: *tcpAddr,
	}

	return wsa, nil
}

func (wst *WsTransport) AllowMulticast() bool {
	return false
}

// Return network, udp, websockets, tcp, ipv4, etc.
func (wst *WsTransport) Network() string {
	return "tcp"
}
