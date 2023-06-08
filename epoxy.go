package main

import "net/http"

func main() {
	websocketHandler := WebsocketHandler()
	receiverHandler := RecieverHandler()

	websocketMux := http.NewServeMux()
	receiverMux := http.NewServeMux()

	websocketMux.HandleFunc("/ws", websocketHandler)
	receiverMux.HandleFunc("/", receiverHandler)

	go func() {
		http.ListenAndServe(":8080", websocketMux)
	}()

	go func() {
		http.ListenAndServe(":8081", receiverMux)
	}()

	select {}

}
