package main

import (
	"flag"

	"github.com/thebluefowl/epoxy"
)

func main() {
	websocketPortPtr := flag.Int("websocket-port", 8136, "websocket port")
	receiverPortPtr := flag.Int("receiver-port", 8137, "receiver port")

	flag.Parse()

	e := epoxy.New(*websocketPortPtr, *receiverPortPtr)
	e.Start()
}
