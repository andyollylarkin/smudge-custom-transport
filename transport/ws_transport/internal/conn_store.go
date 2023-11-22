package internal

import (
	"net"
	"sync"
)

type ConnectionStore struct {
	conns map[string]*WsConnAdapter
	mu    sync.RWMutex
}

func NewConnectionStore() *ConnectionStore {
	return &ConnectionStore{
		conns: make(map[string]*WsConnAdapter),
	}
}

// ConnCacheSet store connection in cache.
func (cs *ConnectionStore) ConnCacheSet(addr net.Addr, conn *WsConnAdapter) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.conns[addr.String()] = conn
}

// ConnCacheRemove remove connection from cache.
func (cs *ConnectionStore) ConnCacheRemove(addr net.Addr) bool {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	delete(cs.conns, addr.String())

	return true
}

// ConnCacheGet get connection from cache.
func (cs *ConnectionStore) ConnCacheGet(addr net.Addr) (*WsConnAdapter, bool) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	conn, ok := cs.conns[addr.String()]
	if !ok {
		return nil, false
	}

	return conn, true
}
