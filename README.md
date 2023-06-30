# Epoxy

Epoxy helps you tunnel requests through the web.  Epoxy was built to evaluate
the feasibility of using asynchronous transport systems like WebSockets to proxy
HTTP traffic.  This projec is not intended to be used in any production systems.
However, you can use this as an alternative to ngrok!  Oh, you can self host the
server in a jiffy too!

## How to Use

```
./epoxy -server <server_address> -target <target_address> -tunnel <tunnel_type> -protocol <protocol_type>
```

Replace <server_address> with the address of the Epoxy server, <target_address> with the address of the target server (default: http://localhost:8000), <tunnel_type> with the desired tunneling channel (e.g., ws for WebSocket), and <protocol_type> with the desired protocol to proxy (e.g., http).

If the <server_address> is not supplied, it will default to the cloud web server.

## Options

- `-server`: The address of the Epoxy server. Defaults to epoxy.everything.sh.
- `-target`: The address of the target server. Defaults to http://localhost:8000.
- `-tunnel`: The tunneling channel to use. Defaults to ws (WebSocket).
- `-protocol`: The target protocol to proxy. Defaults to http.

## Roadmap
- Support other asynchronous channels for tunneling.
- Support for TCP traffic.