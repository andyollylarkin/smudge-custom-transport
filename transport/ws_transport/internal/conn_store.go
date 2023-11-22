package internal

import (
	"fmt"
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
func (cs *ConnectionStore) ConnCacheSet(addr net.Addr, conn *WsConnAdapter) error {
	h, err := extractIpFromAddr(addr)
	if err != nil {
		return fmt.Errorf("cant set conn cache for %s, %w", addr.String(), err)
	}

	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.conns[h] = conn

	return nil
}

// ConnCacheRemove remove connection from cache.
func (cs *ConnectionStore) ConnCacheRemove(addr net.Addr) bool {
	h, err := extractIpFromAddr(addr)
	if err != nil {
		return false
	}

	cs.mu.Lock()
	defer cs.mu.Unlock()

	delete(cs.conns, h)

	return true
}

// ConnCacheGet get connection from cache.
func (cs *ConnectionStore) ConnCacheGet(addr net.Addr) (*WsConnAdapter, bool, error) {
	h, err := extractIpFromAddr(addr)
	if err != nil {
		return nil, false, fmt.Errorf("cant get conn for addr %s from cache, %w", addr.String(), err)
	}

	cs.mu.Lock()
	defer cs.mu.Unlock()

	conn, ok := cs.conns[h]
	if !ok {
		return nil, false, nil
	}

	return conn, true, nil
}

func extractIpFromAddr(addr net.Addr) (string, error) {
	h, _, err := net.SplitHostPort(addr.String())
	if err != nil {
		return "", err
	}

	return h, nil
}
