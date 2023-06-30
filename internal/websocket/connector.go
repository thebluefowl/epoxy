package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/thebluefowl/epoxy/internal/processor"
	"github.com/thebluefowl/epoxy/internal/registry"
	"github.com/thebluefowl/epoxy/internal/tunnel"
)

type Connector struct {
	port     int
	registry *registry.Registry
}

func NewConnector(port int, registry *registry.Registry) tunnel.Connector {
	return &Connector{
		port:     port,
		registry: registry,
	}
}

func (c *Connector) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tConn := NewConnection(conn)
		c.registry.Register(tConn)
		go func(r tunnel.Reader) {
			for {
				b, err := r.ReadAll()
				if err != nil {
					return
				}
				prefix := uint8(0x01)
				processor := processor.GetProcessor(prefix)
				if processor == nil {
					return
				}
				if err := processor.ProcessResponse(b[1:]); err != nil {
					return
				}
			}
		}(tConn)
	})
	return http.ListenAndServe(fmt.Sprintf(":%d", c.port), mux)
}
