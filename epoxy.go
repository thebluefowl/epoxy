package epoxy

import (
	"fmt"
	"net/http"
)

type Epoxy struct {
	WebsocketPort int
	ReceiverPort  int
}

func New(websocketPort, receiverPort int) *Epoxy {
	return &Epoxy{
		WebsocketPort: websocketPort,
		ReceiverPort:  receiverPort,
	}
}

func (e *Epoxy) Start() {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/ws", WebsocketHandler())
		http.ListenAndServe(fmt.Sprintf(":%d", e.WebsocketPort), mux)
	}()
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", RecieverHandler())
		http.ListenAndServe(fmt.Sprintf(":%d", e.ReceiverPort), mux)
	}()

	select {}
}
