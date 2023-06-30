package websocket

import (
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/thebluefowl/epoxy/internal/tunnel"
)

type Client struct {
	u *url.URL
}

func NewClient(address string) tunnel.Client {
	return &Client{
		u: &url.URL{Scheme: "ws", Host: address, Path: "/ws"},
	}
}

func (c *Client) Connect() (tunnel.Connection, error) {
	conn, _, err := websocket.DefaultDialer.Dial(c.u.String(), nil)
	if err != nil {
		return nil, err
	}
	return NewConnection(conn), nil
}
