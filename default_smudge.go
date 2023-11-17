package smudge

import (
	"context"
	"fmt"
	"net"

	"github.com/andyollylarkin/smudge-custom-transport/transport"
)

// GetNodes get all connected nodes.
func GetNodes() []*Node {
	return AllNodes()
}

func RunGossip(ctx context.Context, trns transport.Transport, listenIp string, listenPort int,
	initialNodeAddr string,
) error {
	var ip net.IP

	var err error

	if listenIp == "" {
		ip, err = GetLocalIP()
		if err != nil {
			return fmt.Errorf("Could not get local ip: %w", err)
		}
	} else {
		ip = net.ParseIP(listenIp)
	}

	SetTransport(trns)
	SetLogThreshold(LogInfo)
	SetListenPort(listenPort)
	SetHeartbeatMillis(heartbeatMillis)
	SetListenIP(ip)

	if ip.To4() == nil {
		SetMaxBroadcastBytes(512) // 512 for IPv6
	}

	if initialNodeAddr != "" {
		node, err := CreateNodeByAddress(initialNodeAddr)

		if err == nil {
			AddNode(node)
		} else {
			fmt.Println(err)
		}
	}

	go func() {
		Begin()
	}()

	<-ctx.Done()

	return nil
}
