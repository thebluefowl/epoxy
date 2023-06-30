package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/thebluefowl/epoxy"
	"github.com/thebluefowl/epoxy/internal/proxy"
	"github.com/thebluefowl/epoxy/internal/tunnel"
	"github.com/thebluefowl/epoxy/internal/websocket"
	"golang.org/x/exp/slog"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	server := flag.String("server", "epoxy.everything.sh", "Epoxy server address")
	target := flag.String("target", "http://localhost:8000", "Target address")

	tunnelType := flag.String("tunnel", "ws", "Tunneling channel to use (e.g. ws)")
	protocolType := flag.String("protocol", "http", "Target protocol to proxy (e.g. http)")

	flag.Parse()

	var connection tunnel.Connection
	var err error

	switch {
	case *tunnelType == epoxy.TunnelTypeWebSocket:
		connector := websocket.NewClient(*server)
		connection, err = connector.Connect()
		if err != nil {
			slog.Error("failed to connect to server, err:", slog.Any("err", err))
			os.Exit(1)
		}
	}
	defer connection.Close()

	var proxier proxy.Proxier

	switch {
	case *protocolType == epoxy.ProtocolTypeHTTP:
		proxier = proxy.NewHTTPProxy(*target)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			b, err := connection.ReadAll()
			if err != nil {
				slog.Error("failed to read from connection, err:", slog.Any("err", err))
				return
			}
			if bytes.HasPrefix(b, []byte(epoxy.EpoxyControlID)) {
				id := bytes.Split(b, []byte("="))
				log.Println("New Client Registered: ", fmt.Sprintf("%s.%s", string(id[1]), *server))
				continue
			}
			response, err := proxier.Do(b)
			if err != nil {
				slog.Error("failed to proxy request, err:", slog.Any("err", err))
				return
			}
			// TODO
			response = append([]byte{0x01}, response...)
			if err := connection.Write(response); err != nil {
				slog.Error("failed to write to connection, err:", slog.Any("err", err))
				return
			}
		}
	}()

	for {
		select {
		case <-done:
			slog.Info("finished")
			return
		case <-interrupt:
			fmt.Println("Closing connection")
			err := connection.Write([]byte(epoxy.EpoxyControlID + "close"))
			if err != nil {
				slog.Error("failed to write to connection, err:", err)
				return
			}
			os.Exit(0)
		}
	}
}
