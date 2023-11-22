package internal

import (
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectionStore_ConnCacheSet(t *testing.T) {
	type fields struct {
		conns map[string]*WsConnAdapter
		mu    sync.RWMutex
	}
	type args struct {
		addr net.Addr
		conn *WsConnAdapter
	}
	tests := []struct {
		name         string
		fields       fields
		expectedAddr string
		args         args
	}{
		{
			name: "Set ok",
			fields: fields{
				conns: make(map[string]*WsConnAdapter),
			},
			args: args{
				addr: &WsAddr{
					WsAddrTCP: func() net.TCPAddr {
						a, _ := net.ResolveTCPAddr("tcp", "192.168.1.1:8888")
						return *a
					}(),
				},
				conn: &WsConnAdapter{
					realRemoteAddr: func() *net.TCPAddr {
						a, _ := net.ResolveTCPAddr("tcp", "192.168.1.1:8888")
						return a
					}(),
				},
			},
			expectedAddr: "192.168.1.1:8888",
		},
		{
			name: "Set cant parse addr",
			fields: fields{
				conns: make(map[string]*WsConnAdapter),
			},
			args: args{
				addr: &WsAddr{
					WsAddrTCP: func() net.TCPAddr {
						a := net.TCPAddr{
							IP: net.IPv4(192, 168, 1, 1),
						}
						return a
					}(),
				},
				conn: &WsConnAdapter{
					realRemoteAddr: func() *net.TCPAddr {
						a := net.TCPAddr{
							IP: net.IPv4(192, 168, 1, 1),
						}
						return &a
					}(),
				},
			},
			expectedAddr: "192.168.1.1:0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &ConnectionStore{
				conns: tt.fields.conns,
				mu:    tt.fields.mu,
			}
			cs.ConnCacheSet(tt.args.addr, tt.args.conn)
			require.NotNil(t, tt.args.conn)
			assert.Equal(t, tt.expectedAddr, tt.args.conn.RemoteAddr().String())
		})
	}
}

func TestConnectionStore_ConnCacheGet(t *testing.T) {
	type fields struct {
		conns map[string]*WsConnAdapter
		mu    sync.RWMutex
	}
	type args struct {
		addr                 net.Addr
		conn                 *WsConnAdapter
		setConnectionToCache bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *WsConnAdapter
		want1   bool
	}{
		{
			name: "Get connection ok",
			fields: fields{
				conns: make(map[string]*WsConnAdapter),
			},
			args: args{
				addr: &WsAddr{
					WsAddrTCP: net.TCPAddr{
						IP: net.IPv4(192, 168, 1, 1),
					},
				},
				setConnectionToCache: true,
				conn:                 &WsConnAdapter{},
			},
			want:    &WsConnAdapter{},
			want1:   true,
		},
		{
			name: "Get connection when connection not set",
			fields: fields{
				conns: make(map[string]*WsConnAdapter),
			},
			args: args{
				addr: &WsAddr{
					WsAddrTCP: net.TCPAddr{
						IP: net.IPv4(192, 168, 1, 1),
					},
				},
				setConnectionToCache: false,
				conn:                 &WsConnAdapter{},
			},
			want:    &WsConnAdapter{},
			want1:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &ConnectionStore{
				conns: tt.fields.conns,
			}
			if tt.args.setConnectionToCache {
				cs.ConnCacheSet(tt.args.addr, tt.want)
			}
			_, ok := cs.ConnCacheGet(tt.args.addr)
			assert.Equal(t, tt.want1, ok)
		})
	}
}
