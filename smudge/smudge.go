/*
Copyright 2016 The Smudge Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/andyollylarkin/smudge-custom-transport"
	updtransport "github.com/andyollylarkin/smudge-custom-transport/transport/upd_transport"
)

func main() {
	var nodeAddress string
	var heartbeatMillis int
	var listenPort int
	var listenIp string
	var err error

	flag.StringVar(&nodeAddress, "node", "", "Initial node")
	flag.StringVar(&listenIp, "ip", "", "Listen addr")

	flag.IntVar(&listenPort, "port",
		int(smudge.GetListenPort()),
		"The bind port")

	flag.IntVar(&heartbeatMillis, "hbf",
		int(smudge.GetHeartbeatMillis()),
		"The heartbeat frequency in milliseconds")

	flag.Parse()

	var ip net.IP

	if listenIp == "" {
		ip, err = smudge.GetLocalIP()
		if err != nil {
			log.Fatal("Could not get local ip:", err)
		}
	} else {
		ip = net.ParseIP(listenIp)
	}

	smudge.SetTransport(&updtransport.UDPTransport{})
	smudge.SetLogThreshold(smudge.LogInfo)
	smudge.SetListenPort(listenPort)
	smudge.SetHeartbeatMillis(heartbeatMillis)
	smudge.SetListenIP(ip)

	if ip.To4() == nil {
		smudge.SetMaxBroadcastBytes(512) // 512 for IPv6
	}

	if nodeAddress != "" {
		node, err := smudge.CreateNodeByAddress(nodeAddress)

		if err == nil {
			smudge.AddNode(node)
		} else {
			fmt.Println(err)
		}
	}

	go func() {
		for {

			time.Sleep(time.Second * 5)
			for _, n := range smudge.AllNodes() {
				fmt.Println(n.Address(), n.Status())
			}
		}
	}()

	if err == nil {
		smudge.Begin()
	} else {
		fmt.Println(err)
	}
}
