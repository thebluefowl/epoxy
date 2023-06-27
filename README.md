# Epoxy

Epoxy is a Go project that provides a utility for translating HTTP requests into WebSocket or similar connection protocols. This tool is particularly useful when you need to send HTTP requests over a connection established from the client.

## How it Works

Epoxy acts as a middleware layer that intercepts incoming HTTP requests and translates them into the appropriate WebSocket or similar connection protocol messages. This allows you to seamlessly communicate with a server that expects a different communication protocol, such as WebSocket, while still using standard HTTP requests.