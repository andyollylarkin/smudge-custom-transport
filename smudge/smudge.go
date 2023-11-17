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
	"context"
	"flag"

	"github.com/andyollylarkin/smudge-custom-transport"
)

func main() {
	var nodeAddress string
	var heartbeatMillis int
	var listenPort int
	var listenIp string

	flag.StringVar(&nodeAddress, "node", "", "Initial node")
	flag.StringVar(&listenIp, "ip", "", "Listen addr")

	flag.IntVar(&listenPort, "port",
		int(smudge.GetListenPort()),
		"The bind port")

	flag.IntVar(&heartbeatMillis, "hbf",
		int(smudge.GetHeartbeatMillis()),
		"The heartbeat frequency in milliseconds")

	flag.Parse()

	smudge.RunGossip(context.Background(), nil, listenIp, listenPort, nodeAddress)
}
