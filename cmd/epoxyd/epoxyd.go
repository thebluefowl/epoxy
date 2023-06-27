package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	EpoxyHeaderRequestID = "X-Epoxy-Request-Id"
)

var client = http.DefaultClient

func handleRequest(targetHost string, conn *websocket.Conn, b []byte) {
	rr := bufio.NewReader(bytes.NewReader(b))
	reqIn, err := http.ReadRequest(rr)
	if err != nil {
		fmt.Println(err)
		return
	}

	id := reqIn.Header.Get(EpoxyHeaderRequestID)

	reqOut, err := http.NewRequest(
		reqIn.Method,
		targetHost+reqIn.URL.Path,
		reqIn.Body,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, v := range reqIn.Header {
		reqOut.Header[k] = v
	}

	res, err := client.Do(reqOut)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	res.Header.Set(EpoxyHeaderRequestID, id)

	responseBytes, err := httputil.DumpResponse(res, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, responseBytes); err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	server := flag.String("server", "localhost:8136", "server address")
	port := flag.String("port", "8136", "local http port")

	flag.Parse()

	targetHost := fmt.Sprintf("http://localhost:%s", *port)

	u := url.URL{Scheme: "ws", Host: *server, Path: "/ws"}
	headers := make(http.Header)
	headers.Set("X-Epoxy-Source", uuid.New().String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, b, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			handleRequest(targetHost, conn, b)
		}
	}()

	for {
		select {
		case <-done:
			fmt.Println("Finished")
			return
		case <-interrupt:
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			<-done
			return
		}
	}
}
