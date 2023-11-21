package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseURIToHostPort(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name     string
		args     args
		wantIp   string
		wantPort int
		wantErr  bool
	}{
		{
			name: "Correct url",
			args: args{
				uri: "http://192.168.10.10:80/",
			},
			wantIp:   "192.168.10.10",
			wantPort: 80,
			wantErr:  false,
		},
		{
			name: "Correct url without port",
			args: args{
				uri: "http://192.168.10.10/",
			},
			wantIp:   "192.168.10.10",
			wantPort: 80,
			wantErr:  false,
		},
		{
			name: "Correct url not default port",
			args: args{
				uri: "http://192.168.10.10:8080/",
			},
			wantIp:   "192.168.10.10",
			wantPort: 8080,
			wantErr:  false,
		},
		{
			name: "Incorrect url",
			args: args{
				uri: "http://192.16",
			},
			wantIp:   "",
			wantPort: 0,
			wantErr:  true,
		},
		{
			name: "Incorrect url (2). Without scheme",
			args: args{
				uri: "192.168.1.1:80",
			},
			wantIp:   "",
			wantPort: 0,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIp, gotPort, err := ParseURIToHostPort(tt.args.uri)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantIp, gotIp)
			assert.Equal(t, tt.wantPort, gotPort)
		})
	}
}
