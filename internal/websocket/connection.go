package websocket

import (
	"github.com/gorilla/websocket"
)

type Connection struct {
	conn *websocket.Conn
}

func NewConnection(conn *websocket.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) Write(b []byte) error {
	return c.conn.WriteMessage(websocket.BinaryMessage, b)
}

func (c *Connection) ReadAll() ([]byte, error) {
	_, b, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

func (c *Connection) End() error {
	return c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
