package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/thebluefowl/epoxy/internal/processor"
	"github.com/thebluefowl/epoxy/internal/registry"
	"github.com/thebluefowl/epoxy/internal/websocket"
)

func main() {
	websocketPort := flag.Int("ws-port", 8136, "websocket port")
	httpPort := flag.Int("http-port", 8137, "http port")
	httpTimeout := flag.Duration("http-timeout", 5*time.Second, "http timeout")

	flag.Parse()

	registry := registry.NewRegistry()

	connector := websocket.NewConnector(*websocketPort, registry)

	go func() {
		err := connector.Start()
		if err != nil {
			fmt.Println("ERR: failed to start websocket connector, err:", err)
			os.Exit(1)
		}
	}()

	httpProcessor := processor.NewHTTPProcessor(registry)
	processor.RegisterProcessor(uint8(0x01), httpProcessor)

	go func() {
		err := httpProcessor.Start(&processor.HTTPOpts{
			Port:    *httpPort,
			Timeout: *httpTimeout,
		})
		if err != nil {
			fmt.Println("ERR: failed to start http processor, err:", err)
			os.Exit(1)
		}
	}()

	select {}

}
