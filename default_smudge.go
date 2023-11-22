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

func ThisHost() *Node {
	return thisHost
}

func RunGossip(ctx context.Context, trns transport.Transport, listenIp string, listenPort int,
	initialNodeAddr string, logger Logger, logLvl LogLevel,
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

	if logger != nil {
		SetLogger(logger)
	} else {
		SetLogThreshold(LogAll)
	}

	SetLogThreshold(logLvl)
	SetListenPort(listenPort)
	SetHeartbeatMillis(heartbeatMillis)
	SetListenIP(ip)

	// Redefine address to iface ipv4 address
	if listenIp == "" {
		if localIpv4Addr, err := TryGetLocalIPv4(); err != nil || localIpv4Addr.To4() == nil {
			return fmt.Errorf("cant redefine listen interface to IPv4 address. Error: %w", err)
		} else {
			logger.Logf(LogInfo, "Listen IPv4 address: %s", localIpv4Addr.String())
			SetListenIP(localIpv4Addr)
		}
	}

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
