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
	"github.com/andyollylarkin/smudge-custom-transport/pkg/utils"
	"github.com/andyollylarkin/smudge-custom-transport/transport"
	"github.com/andyollylarkin/smudge-custom-transport/transport/ws_transport/internal"
	"github.com/gorilla/websocket"
	lru "github.com/hashicorp/golang-lru"
)

const (
	MaxLRUCacheItems int = 100
)

var (
	upgrader websocket.Upgrader
)

type WsTransport struct {
	cache              *lru.Cache
	wg                 sync.WaitGroup
	remoteWsServerPort int
	connChan           chan transport.GenericConn
	logger             smudge.Logger
}

// connCacheSet store connection in LRU cache.
func (wst *WsTransport) connCacheSet(addr net.Addr, conn *internal.WsConnAdapter) (bool, error) {
	h, _, err := net.SplitHostPort(addr.String())
	if err != nil {
		return false, fmt.Errorf("cant set conn cache for %s, %w", addr.String(), err)
	}

	return wst.cache.Add(h, conn), nil
}

// connCacheRemove remove connection from LRU cache.
func (wst *WsTransport) connCacheRemove(addr net.Addr) bool {
	h, _, err := net.SplitHostPort(addr.String())
	if err != nil {
		return false
	}

	return wst.cache.Remove(h)
}

// connCacheGet get connection from LRU cache.
func (wst *WsTransport) connCacheGet(addr net.Addr) (*internal.WsConnAdapter, bool, error) {
	h, _, err := net.SplitHostPort(addr.String())
	if err != nil {
		return nil, false, fmt.Errorf("cant get conn for addr %s from cache, %w", addr.String(), err)
	}

	conn, ok := wst.cache.Get(h)
	if !ok {
		wst.logger.Logf(smudge.LogDebug, "connection get cache miss for %s", addr.String())

		return nil, false, nil
	}

	wsConn, ok := conn.(*internal.WsConnAdapter)
	if !ok {
		return nil, false, fmt.Errorf("cant get conn for addr %s from cache. Conn type isn't WsConnAdapter", addr.String())
	}

	return wsConn, true, nil
}

func NewWsTransport(logger smudge.Logger, remoteWsServerPort int) (*WsTransport, error) {
	cache, err := lru.New(MaxLRUCacheItems)
	if err != nil {
		return nil, fmt.Errorf("cant create connections cache: %w", err)
	}

	t := new(WsTransport)

	t.logger = logger
	t.remoteWsServerPort = remoteWsServerPort
	t.connChan = make(chan transport.GenericConn)

	t.cache = cache

	return t, nil
}

// UpgageWebsocket upgrade http request to websocket connection. Pass it to web server handler.
func (wst *WsTransport) UpgageWebsocket(w http.ResponseWriter, r *http.Request) error {
	wsconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("cant upgrade websocket connection: %w", err)
	}

	_, ok, err := wst.connCacheGet(wsconn.RemoteAddr())
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	adapter := &internal.WsConnAdapter{
		SocketConn: wsconn,
	}

	_, err = wst.connCacheSet(wsconn.RemoteAddr(), adapter)
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
		conn, ok, err := wst.connCacheGet(addr)
		if err != nil || !ok {
			continue
		}

		conn.ActuallyClose()

		wst.connCacheRemove(addr)

		wst.logger.Logf(smudge.LogDebug, "Actually close %s", conn.RemoteAddr().String())
	}
}

func (wst *WsTransport) Dial(ctx context.Context, laddr transport.SockAddr,
	raddr transport.SockAddr,
) (transport.GenericConn, error) {
	c, ok, err := wst.connCacheGet(raddr)
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

	ip, port, err := utils.ParseURIToHostPort(raddr.String())
	if err != nil {
		return nil, err
	}

	remoteAddr := net.JoinHostPort(ip, strconv.Itoa(port))

	url := url.URL{
		Scheme: "ws",
		Host:   remoteAddr,
		Path:   WebsocketRoutePath,
	}

	wsconn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, err
	}

	adapter := &internal.WsConnAdapter{
		SocketConn: wsconn,
	}

	_, err = wst.connCacheSet(adapter.RemoteAddr(), adapter)
	if err != nil {
		return nil, err
	}

	wst.connChan <- adapter

	return adapter, nil
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
