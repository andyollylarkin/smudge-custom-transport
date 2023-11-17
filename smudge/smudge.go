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
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/andyollylarkin/smudge-custom-transport"
	"github.com/andyollylarkin/smudge-custom-transport/pkg/logger"
	wstransport "github.com/andyollylarkin/smudge-custom-transport/transport/ws_transport"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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

	l := logger.NewLogrusLogger(logrus.New(), logrus.DebugLevel)

	t, err := wstransport.NewWsTransport(l)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc(wstransport.WebsocketRoutePath, func(w http.ResponseWriter, r *http.Request) {
		err := t.UpgageWebsocket(w, r)
		if err != nil {
			log.Println(err)
		}
	})

	go smudge.RunGossip(context.Background(), t, listenIp, listenPort, nodeAddress, l, smudge.LogAll)
	// smudge.RunGossip(context.Background(), nil, listenIp, listenPort, nodeAddress)

	http.ListenAndServe(net.JoinHostPort(listenIp, strconv.Itoa(listenPort)), r)
}
