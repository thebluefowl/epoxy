package main

import (
	"bufio"
	"bytes"
	"errors"
	"net/http"
	"net/http/httputil"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	ErrInvalidWritePayload = errors.New("invalid write payload")
)

type websocketConn struct {
	conn *websocket.Conn
	lock *sync.Mutex
}

func newWebsocketConn(conn *websocket.Conn) ReadWriter {
	return &websocketConn{
		conn: conn,
		lock: &sync.Mutex{},
	}
}

func (c *websocketConn) Write(request *http.Request) error {
	payload, err := httputil.DumpRequest(request, true)
	if err != nil {
		return errors.Join(err, ErrInvalidWritePayload)
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.conn.WriteMessage(websocket.BinaryMessage, payload); err != nil {
		return err
	}
	return nil
}

func (c *websocketConn) Read() (*http.Response, error) {
	_, payload, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	response, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(payload)), nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *websocketConn) Close() {
	c.conn.Close()
}

func WebsocketHandler() http.HandlerFunc {
	upgrader := websocket.Upgrader{}
	return func(w http.ResponseWriter, r *http.Request) {
		target := r.Header.Get("X-Epoxy-Source")
		if target == "" {
			http.Error(w, "Missing X-Target header", http.StatusBadRequest)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
			return
		}
		clientRegistry.Register(target, newWebsocketConn(conn))
		clientRegistry.Listen(target)
	}
}
